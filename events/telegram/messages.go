// Package telegram ...
package telegram

const msgHelp = `Я умею сохранять твои грязные URLишки. Также я могу предложить тебе порыться в них.

Чтобы сохранить страницу, просто пришли мне ссылку на нее.

Чтобы получить случайную страницу из списка, отправь мне команду /rnd.
Бро будь осторожен! После этого страница будет удалена из списка сохраненных!`

const msgHello = "Превед! \n\n" + msgHelp

const (
	msgUnknownCommand = "Твоя команда 💩"
	msgNoSavedPages   = "У тебя нет сохраненных URLишек 👀"
	msgSaved          = "Ауф! 👍"
	msgAlredyExists   = "У тебя уже есть это дерьмо 🤦‍♂"
)
