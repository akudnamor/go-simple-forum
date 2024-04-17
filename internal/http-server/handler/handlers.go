package handler

import (
	"go-simple-forum/internal/lib"
	"go-simple-forum/internal/lib/logger/sl"
	"go-simple-forum/internal/storage"
	"go-simple-forum/internal/viewmodel"
	"html/template"

	"log/slog"
	"net/http"
	"strconv"
)

func IndexPage(log *slog.Logger, st *storage.Storage, t *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vml := viewmodel.LoggedIn{}
		vml.Status = "False"
		vml.User = storage.User{}
		c, err := r.Cookie("sessionId")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Info("no cookie")

			} else {
				log.Info("some other problem")
			}

			//случай когда нету куки:
			err = t.ExecuteTemplate(w, "index", vml)
			if err != nil {
				log.Error("failed to execute template", sl.Err(err))
			}
			return
		}
		log.Info("/index -> sessionId:", c.Value)

		// случай если куки есть:
		sessionId, err := strconv.Atoi(c.Value)
		if err != nil {
			log.Error("cant atoi", sl.Err(err))
		}

		session, err := st.GetSessionById(sessionId)
		if err != nil {
			log.Error("cant get session by id", sl.Err(err))
		}

		user, err := st.GetUserByUserId(session.UserId)
		if err != nil {
			log.Error("cant get user by userid", sl.Err(err))
		}

		vml.Status = "True"
		vml.User = user
		log.Debug("info in vml", vml)
		err = t.ExecuteTemplate(w, "index", vml)
		if err != nil {
			log.Error("failed to execute template", sl.Err(err))
		}

		return
	}
}

func SignUpPage(log *slog.Logger, st *storage.Storage, t *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("sessionId")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Info("no cookie")

			} else {
				log.Info("some other problem")
			}

			//случай когда нету куки + ретурн в конце
			err = t.ExecuteTemplate(w, "signup", nil)
			if err != nil {
				log.Error("failed to execute template", sl.Err(err))
			}
			return
		}
		log.Info("/signup -> sessionId:", c.Value)

		// случай если куки есть:
		sessionId, err := strconv.Atoi(c.Value)
		if err != nil {
			log.Error("cant atoi", sl.Err(err))
		}

		session, err := st.GetSessionById(sessionId)
		if err != nil {
			log.Error("cant get session by id", sl.Err(err))
		}

		user, err := st.GetUserByUserId(session.UserId)
		if err != nil {
			log.Error("cant get user by userid", sl.Err(err))
		}

		log.Info("info about this user:", user)

		err = t.ExecuteTemplate(w, "signup", user)
		if err != nil {
			log.Error("failed to execute template", sl.Err(err))
		}

		return
	}
}

func AuthPage(log *slog.Logger, st *storage.Storage, t *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("sessionId")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Info("no cookie")

			} else {
				log.Info("some other problem")
			}

			//случай когда нету куки + ретурн в конце
			err = t.ExecuteTemplate(w, "auth", nil)
			if err != nil {
				log.Error("failed to execute template", sl.Err(err))
			}
			return
		}
		log.Info("/auth -> sessionId:", c.Value)

		// случай если куки есть:
		sessionId, err := strconv.Atoi(c.Value)
		if err != nil {
			log.Error("cant atoi", sl.Err(err))
		}

		session, err := st.GetSessionById(sessionId)
		if err != nil {
			log.Error("cant get session by id", sl.Err(err))
		}

		user, err := st.GetUserByUserId(session.UserId)
		if err != nil {
			log.Error("cant get user by userid", sl.Err(err))
		}

		log.Info("info about this user:", user)

		err = t.ExecuteTemplate(w, "auth", user)
		if err != nil {
			log.Error("failed to execute template", sl.Err(err))
		}

		return
	}
}

func Auth(log *slog.Logger, st *storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		login := r.FormValue("login")
		password := r.FormValue("password")

		log.Debug("get from form:", login, password)

		users, err := st.GetUsers()
		if err != nil {
			log.Error("cant get users", sl.Err(err))
		}
		log.Debug("users from db:", users)
		for _, user := range users {

			//log.Debug("", idx, "user from db:", user.Login, user.Password)

			if login == user.Login && password == user.Password {
				log.Debug("login success")
				sessionId, err := st.AddSession(user.UserId)
				if err != nil {
					log.Error("Cant add session", sl.Err(err))
				}
				cookie := http.Cookie{
					Name:  "sessionId",
					Value: strconv.Itoa(sessionId),
					Path:  "/",
				}
				http.SetCookie(w, &cookie)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		http.Redirect(w, r, "/auth", http.StatusSeeOther)

	}
}

func SignUp(log *slog.Logger, st *storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		login := r.FormValue("login")
		password := r.FormValue("password")
		userId, err := lib.GenerateUserId()
		if err != nil {
			log.Error("Cant generate UserId", sl.Err(err))
		}
		err = st.AddUser(login, password, userId)
		if err != nil {
			log.Error("cant add new user", sl.Err(err))
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Logout(log *slog.Logger, st *storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("sessionId")
		if err != nil {
			log.Error("Cant parse cookie session", sl.Err(err))
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		sessionId := c.Value
		err = st.DeleteSessionBySessionId(sessionId)
		if err != nil {
			log.Error("cant delete session", sl.Err(err))
		}

		cookie := http.Cookie{
			Name:   "sessionId",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	}
}
