package middleware

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/Nidal-Bakir/first_go/pkg/tracker"
	"github.com/Nidal-Bakir/first_go/pkg/user"
	"github.com/google/uuid"
)

type middlewareCtxKey int

const ()

func extractUserFromRequestCookie(req *http.Request) (user.UserModel, error) {
	cookie, err := req.Cookie("user")

	if err != nil {
		return user.UserModel{}, err
	}
	userName := cookie.Value
	return user.UserModel{Name: userName}, nil
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userModel, err := extractUserFromRequestCookie(r)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		// we should check with the DB about this user first but we will skip this fo now!!
		r = r.WithContext(user.ContextWithUser(r.Context(), userModel))
		h.ServeHTTP(w, r)
	})
}

func WithAuthUserHandlerFunc(handler func(w http.ResponseWriter, r *http.Request, user user.UserModel)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := user.UserFromContext(r.Context())
		if !ok {
			log.Println("userFromContext returned ok=false")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		handler(w, r, user)
	}
}

func RequestUUIDMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var uuidVal uuid.UUID

		uuidStr := r.Header.Get("X-Request-UUID")
		if uuidStr == "" {
			uuidVal = uuid.New()
		} else {
			if u, err := uuid.Parse(uuidStr); err == nil {
				uuidVal = u
			} else {
				log.Println(err)
				uuidVal = uuid.New()
			}
		}

		ctx = tracker.ContextWithReqUUID(ctx, uuidVal)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ClientIPMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ip := clientIP(r); ip != "" {
			ctx = tracker.ContextWithClientIP(ctx, ip)
		}
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func clientIP(r *http.Request) string {
	var ip string

	if trueIP := r.Header.Get("True-Client-IP"); trueIP != "" {
		ip = trueIP
	} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		ip = realIP
	} else if xffIPs := r.Header.Get("X-Forwarded-For"); xffIPs != "" {
		i := strings.Index(xffIPs, ",")
		if i == -1 {
			i = len(xffIPs)
		}
		ip = strings.TrimSpace(xffIPs[:i])
	}

	if isValidIP(ip) {
		return ip
	}

	return ""
}

func LastIPMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		xffHeaders := r.Header.Values("X-Forwarded-For")
		if ip := lastIP(xffHeaders); ip != "" {
			ctx = tracker.ContextWithLastIP(ctx, ip)
		}

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func lastIP(xffHeaders []string) string {
	l := len(xffHeaders)
	if l == 0 {
		return ""
	}

	var lastIP string

	xffHeader := joinXFFHeaders(xffHeaders)

	lastIPIndex := strings.LastIndex(xffHeader, ",")
	if lastIPIndex == -1 {
		lastIP = xffHeader
	} else {
		lastIP = xffHeader[lastIPIndex+1:]
	}

	lastIP = strings.TrimSpace(lastIP)
	if isValidIP(lastIP) {
		return lastIP
	}

	return ""
}
func isValidIP(ip string) bool {
	return ip != "" && net.ParseIP(ip) != nil
}

func joinXFFHeaders(xffHeaders []string) string {
	return strings.Join(xffHeaders, ", ")
}
