package bot

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/aaltgod/telegram-bot/pkg/ipstack"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	InternalError = "Внутренняя ошибка. Повторите запрос позже"
)

var (
	reIPv4     = regexp.MustCompile(`(\b25[0-5]|\b2[0-4][0-9]|\b[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	reIPv6     = regexp.MustCompile(`(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`)
	operations = make(map[int]operationType)
)

type (
	operation     string
	operationType int
)

const (
	userMenuOperationType operationType = iota
	adminMenuOperationType

	newAdminOperationType
	deleteAdminOperationType
	userStatOperationType
	distributionOperationType
	checkIPOperationType
	requestHistoryOperationType
)

var (
	startOperation  = operation("start")
	menuOperation   = operation("menu")
	cancelOperation = operation("cancel")

	newAdminOperation       = operation("Добавить админа")
	deleteAdminOperation    = operation("Удалить админа")
	userStatOperation       = operation("Статистика пользователя")
	distributionOperation   = operation("Рассылка")
	checkIPOperation        = operation("Проверить IP")
	requestHistoryOperation = operation("История запросов")
)

func (b *Bot) handleCommand(msg *tgbotapi.Message) error {
	b.logger.Infof(
		"handlecommand [%s] with username: %s; id: %d\n",
		msg.Command(), msg.From.UserName, msg.From.ID,
	)

	var msgConfig tgbotapi.MessageConfig

	switch operation(msg.Command()) {
	case startOperation:
		createUser := &CreateUser{
			Name:    msg.From.UserName,
			ID:      int64(msg.From.ID),
			IsAdmin: false,
		}

		if err := b.api.InsertUser(createUser); err != nil {
			return err
		}

		msgConfig = tgbotapi.NewMessage(msg.Chat.ID, "Привет")

	case menuOperation:
		user, err := b.api.GetUser(int64(msg.From.ID))
		if err != nil {
			return err
		}

		var (
			rowUp   []tgbotapi.KeyboardButton
			rowDown []tgbotapi.KeyboardButton
		)

		if user.IsAdmin {
			rowUp = append(rowUp, tgbotapi.NewKeyboardButton(string(distributionOperation)))
			rowUp = append(rowUp, tgbotapi.NewKeyboardButton(string(userStatOperation)))
			rowDown = append(rowDown, tgbotapi.NewKeyboardButton(string(newAdminOperation)))
			rowDown = append(rowDown, tgbotapi.NewKeyboardButton(string(deleteAdminOperation)))

			operations[msg.From.ID] = adminMenuOperationType
		} else {
			rowUp = append(rowUp, tgbotapi.NewKeyboardButton(string(checkIPOperation)))
			rowDown = append(rowDown, tgbotapi.NewKeyboardButton(string(requestHistoryOperation)))

			operations[msg.From.ID] = userMenuOperationType
		}

		msgConfig = tgbotapi.NewMessage(msg.Chat.ID, "Выберите команду:")
		msgConfig.ReplyMarkup = tgbotapi.NewReplyKeyboard(rowUp, rowDown)

	case cancelOperation:
		if opType, ok := operations[msg.From.ID]; !ok {
			msgConfig = tgbotapi.NewMessage(msg.Chat.ID, "Нет активной команды, которую можно отменить")
		} else {
			msgConfig = tgbotapi.NewMessage(msg.Chat.ID, "Последняя команда была отменена")
			if opType == userMenuOperationType || opType == adminMenuOperationType {
				msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			}
			delete(operations, msg.From.ID)
		}

	default:
		return nil
	}

	b.botApi.Send(msgConfig)
	return nil
}

func (b *Bot) handleText(msg *tgbotapi.Message) error {
	b.logger.Infof(
		"handletext [%s] with username: %s; id: %d\n",
		msg.Text, msg.From.UserName, msg.From.ID,
	)

	opType, ok := operations[msg.From.ID]
	if !ok {
		b.logger.Infof("user hasn't an operation with id: %d", msg.From.ID)
		return nil
	}

	var msgConfig tgbotapi.MessageConfig

	switch opType {
	case userMenuOperationType:
		switch operation(msg.Text) {
		case checkIPOperation:
			operations[msg.From.ID] = checkIPOperationType
			msgConfig = tgbotapi.NewMessage(msg.Chat.ID, "Введите IP")

		case requestHistoryOperation:
			requests, err := b.api.GetAllRequestsByID(int64(msg.From.ID))
			if err != nil {
				return err
			}

			var history strings.Builder

			for _, r := range requests {
				history.WriteString(r.Response)
				history.WriteRune('\n')
				history.WriteRune('\n')
			}

			msgConfig = tgbotapi.NewMessage(msg.Chat.ID, history.String())
			delete(operations, msg.From.ID)
		}
		msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	case adminMenuOperationType:
		switch operation(msg.Text) {
		case newAdminOperation:
			operations[msg.From.ID] = newAdminOperationType
			msgConfig = tgbotapi.NewMessage(msg.Chat.ID, "Введите ID пользователя")

		case deleteAdminOperation:
			operations[msg.From.ID] = deleteAdminOperationType
			msgConfig = tgbotapi.NewMessage(msg.Chat.ID, "Введите ID админа")

		case distributionOperation:
			operations[msg.From.ID] = distributionOperationType
			msgConfig = tgbotapi.NewMessage(msg.Chat.ID, "Отправьте сообщение для рассылки")

		case userStatOperation:
			operations[msg.From.ID] = userStatOperationType
			msgConfig = tgbotapi.NewMessage(msg.Chat.ID, "Введите ID пользователя")
		}
		msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	case checkIPOperationType:
		if reIPv4.MatchString(msg.Text) || reIPv6.MatchString(msg.Text) {
			info, err := ipstack.GetInfo(msg.Text)
			if err != nil {
				return err
			}

			r := &Request{
				IP:       msg.Text,
				Response: info,
			}

			if err := b.api.AppendRequest(int64(msg.From.ID), r); err != nil {
				return err
			}

			msgConfig = tgbotapi.NewMessage(msg.Chat.ID, info)
			delete(operations, msg.From.ID)
		} else {
			msgConfig = tgbotapi.NewMessage(msg.Chat.ID, "Неправильный формат")
		}

	case newAdminOperationType:
		updateUser := &UpdateUser{
			IsAdmin: true,
		}

		id, err := strconv.Atoi(msg.Text)
		if err != nil {
			return err
		}

		if err := b.api.UpdateUser(int64(id), updateUser); err != nil {
			return err
		}
		msgConfig = tgbotapi.NewMessage(msg.Chat.ID, "Пользователь с ID "+msg.Text+" теперь админ")
		delete(operations, msg.From.ID)

	case deleteAdminOperationType:
		updateUser := &UpdateUser{
			IsAdmin: false,
		}

		id, err := strconv.Atoi(msg.Text)
		if err != nil {
			return err
		}

		if err := b.api.UpdateUser(int64(id), updateUser); err != nil {
			return err
		}
		msgConfig = tgbotapi.NewMessage(msg.Chat.ID, "Пользователь с ID "+msg.Text+" больше не админ")
		delete(operations, msg.From.ID)

	case userStatOperationType:
		id, err := strconv.Atoi(msg.Text)
		if err != nil {
			return err
		}

		requests, err := b.api.GetAllRequestsByID(int64(id))
		if err != nil {
			return err
		}

		var history strings.Builder

		for _, r := range requests {
			history.WriteString(r.IP)
			history.WriteRune('\n')
		}

		msgConfig = tgbotapi.NewMessage(msg.Chat.ID, history.String())
		delete(operations, msg.From.ID)

	case distributionOperationType:
		users, err := b.api.GetUsers()
		if err != nil {
			return err
		}

		for _, u := range users {
			if !u.IsAdmin {
				msgConfig = tgbotapi.NewMessage(u.ID, msg.Text)
				b.botApi.Send(msgConfig)
			}
		}

		msgConfig = tgbotapi.NewMessage(msg.Chat.ID, "Произведена рассылка")
		delete(operations, msg.From.ID)
	default:
		return nil
	}

	b.botApi.Send(msgConfig)
	return nil
}

func (b *Bot) handleError(msg *tgbotapi.Message) {
	b.botApi.Send(tgbotapi.NewMessage(msg.Chat.ID, InternalError))
}
