package main

import (
	"context"
	"fmt"
	"github.com/argo0n/TTImmersn_assignment_1/rpc-server/db"
	"github.com/argo0n/TTImmersn_assignment_1/rpc-server/kitex_gen/rpc"
	"math/rand"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

var database, connErr = db.CreateDB("db", "tiktok_chat", "admin", "p@ssw0rd")

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	resp := rpc.NewSendResponse()

	if connErr != nil {
		resp.Code, resp.Msg = 500, "Couldn't connect to database"
		return resp, conn
	}

	_, err := database.ExecInsert(
		"INSERT INTO messages(chat, sender, text, send_time) VALUES(?, ?, ?, ?)",
		req.Message.Chat, req.Message.Text, req.Message.Sender, req.Message.SendTime,
	)
	if err != nil {
		resp.Code, resp.Msg = 500, "Failed to add message to database"
		return resp, err
	}
	println("Success")
	resp.Code, resp.Msg = 0, "success"
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()

	if connErr != nil {
		resp.Code, resp.Msg = 500, "Couldn't connect to database"
		return resp, connErr
	}

	cursor, limit, reverse := req.Cursor, req.Limit, req.Reverse
	if cursor < 0 {
		resp.Code, resp.Msg = 400, "Cursor cannot be less than 0"
		return resp, nil
	}

	order := "ASC"
	if *reverse {
		order = "DESC"
	}

	query := fmt.Sprintf("SELECT chat, text, sender, datetime FROM messages WHERE id > ? ORDER BY id %s LIMIT ?", order)
	rows, err := database.ExecSelectMany(query, cursor, limit)
	if err != nil {
		resp.Code, resp.Msg = 500, "Couldn't execute database query"
		return resp, err
	}

	messages := make([]*rpc.Message, 0)
	for rows.Next() {
		var msg rpc.Message
		err := rows.Scan(&msg.Id, &msg.Chat, &msg.Text, &msg.Sender, &msg.SendTime)
		if err != nil {
			resp.Code, resp.Msg = 500, "Couldn't read data from the database"
			return resp, err
		}
		messages = append(messages, &msg)
	}
	if len(messages) > limit {
		resp.HasMore = true
		resp.NextCursor = messages[limit].SendTime // Set next_cursor as send_time of last fetched message
		messages = messages[:limit]                // Exclude last message from results returned to client
	} else {
		resp.HasMore = false
	}

	resp.Messages = messages
	resp.Code, resp.Msg = areYouLucky()
	return resp, nil
}

func areYouLucky() (int32, string) {
	if rand.Int31n(2) == 1 {
		return 0, "success"
	} else {
		return 500, "oops"
	}
}
