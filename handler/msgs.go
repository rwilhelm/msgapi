package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"git.sr.ht/~rxw/msgapi/db"
	"git.sr.ht/~rxw/msgapi/models"
)

var msgIDKey = "msgID"

func msgs(router chi.Router) {
	router.Get("/", getAllMsgs)
	router.Post("/", createMsg)

	router.Route("/{msgId}", func(router chi.Router) {
		router.Use(MsgContext)
		router.Get("/", getMsg)
		router.Put("/", updateMsg)
		router.Delete("/", deleteMsg)
	})
}

func MsgContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msgId := chi.URLParam(r, "msgId")
		if msgId == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("msg ID is required")))
			return
		}
		id, err := strconv.Atoi(msgId)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid msg ID")))
		}
		ctx := context.WithValue(r.Context(), msgIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAllMsgs(w http.ResponseWriter, r *http.Request) {
	msgs, err := dbInstance.GetAllMsgs()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, msgs); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func createMsg(w http.ResponseWriter, r *http.Request) {
	msg := &models.Msg{}
	if err := render.Bind(r, msg); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := dbInstance.AddMsg(msg); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, msg); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getMsg(w http.ResponseWriter, r *http.Request) {
	msgID := r.Context().Value(msgIDKey).(int)
	msg, err := dbInstance.GetMsgById(msgID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &msg); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func deleteMsg(w http.ResponseWriter, r *http.Request) {
	msgId := r.Context().Value(msgIDKey).(int)
	err := dbInstance.DeleteMsg(msgId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
}

func updateMsg(w http.ResponseWriter, r *http.Request) {
	msgId := r.Context().Value(msgIDKey).(int)
	msgData := models.Msg{}
	if err := render.Bind(r, &msgData); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	msg, err := dbInstance.UpdateMsg(msgId, msgData)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &msg); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
