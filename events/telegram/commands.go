// Package telegram ...
package telegram

import (
	"context"
	"errors"
	"log"
	"net"
	"net/url"
	"strings"

	"github.com/yellowpuki/tg-bot/clients/telegram"
	"github.com/yellowpuki/tg-bot/lib/e"
	"github.com/yellowpuki/tg-bot/storage"
	"golang.org/x/exp/slog"
)

const (
	LastCmd  = "/last"
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

// doCmd ...
func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	if isAddCmd(text) {
		return p.savePage(context.Background(), chatID, text, username)
	}

	switch text {
	case LastCmd:
		return p.sendLast(context.Background(), chatID, username)
	case RndCmd:
		return p.sendRandom(context.Background(), chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

// savePage ...
func (p *Processor) savePage(ctx context.Context, chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	sendMsg := NewMessageSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return err
	}
	if isExists {
		slog.Info("savePage: url exists", slog.String("URL", pageURL))
		return sendMsg(msgAlredyExists)
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

// sendLast ...
func (p *Processor) sendLast(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random page", err) }()

	page, err := p.storage.PickLast(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return nil
}

// sendRandom ...
func (p *Processor) sendRandom(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random page", err) }()

	page, err := p.storage.PickRandom(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(ctx, page)
}

// sendHelp ...
func (p *Processor) sendHelp(chatID int) error {

	return p.tg.SendMessage(chatID, msgHelp)
}

// sendHello ...
func (p *Processor) sendHello(chatID int) error {

	return p.tg.SendMessage(chatID, msgHello)
}

// NewMessageSender ...
func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

// isAddCmd ...
func isAddCmd(text string) bool {
	return isURL(text)
}

// isURL ...
func isURL(text string) bool {
	u, err := url.Parse(text)
	if err != nil {
		slog.Info(err.Error())
		return false
	}

	addr := net.ParseIP(u.Host)

	if addr == nil {
		slog.Info("url-info", "host", u.Host)

		return strings.Contains(u.Host, ".")
	}

	return true
}
