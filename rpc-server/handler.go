package main

import (
	"context"
	"github.com/argo0n/TTImmersn_assignment_1/rpc-server/db"
	"github.com/argo0n/TTImmersn_assignment_1/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

var database, connErr = db.CreateDB("db", "tiktok_chat", "admin", "p@ssw0rd")

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	resp := rpc.NewSendResponse()

	if connErr != nil {
		resp.Code, resp.Msg = 500, "Couldn't connect to database"
		return resp, connErr
	}

	_, err := database.ExecInsert(
		"INSERT INTO messages(chat, sender, text, send_time) VALUES(?, ?, ?, ?)",
		req.Message.Chat, req.Message.Sender, req.Message.Text, req.Message.SendTime,
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
	println("Requested pull")

	if connErr != nil {
		resp.Code, resp.Msg = 500, "Couldn't connect to database"
		return resp, connErr
	}

	cursor, limit, reverse := req.Cursor, req.Limit, req.Reverse
	if cursor < 0 {
		resp.Code, resp.Msg = 400, "Cursor cannot be less than 0"
		return resp, nil
	}

	query := ""
	if *reverse {
		query = "SELECT id, chat, text, sender, send_time AS sendtime FROM messages WHERE id < ? ORDER BY id DESC LIMIT ?"
	} else {
		query = "SELECT id, chat, text, sender, send_time AS sendtime FROM messages WHERE id > ? ORDER BY id LIMIT ?"
	}
	rows, err := database.ExecSelectMany(query, cursor, limit+1)
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
	if len(messages) > int(limit) {
		hasMore := true
		resp.HasMore = &hasMore
		if *reverse {
			nextCursor := cursor - int64(limit)
			if nextCursor < 0 {
				nextCursor = 0
			}
			resp.NextCursor = &nextCursor
		} else {
			nextCursor := cursor + int64(limit)
			resp.NextCursor = &nextCursor
		}

		messages = messages[:int(limit)] // Exclude last message from results returned to client
	} else {
		hasMore := false
		resp.HasMore = &hasMore
	}
	resp.Messages, resp.Code, resp.Msg = messages, 0, "Success"
	return resp, nil
}
