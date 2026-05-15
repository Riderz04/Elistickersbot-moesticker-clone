package core

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

func sendStartMessage(c tele.Context) error {
	// Enviar la imagen con el texto como caption
	photo := &tele.Photo{
		File: tele.FromURL("https://files.catbox.moe/inn3u2.jpg"),
		Caption: `¡Hola! soy <b>Eli♡Bot Stickers🫟</b>, un clon modificado de <a href="https://github.com/star-39/moe-sticker-bot">Moe Sticker Bot</a>. Usa los siguientes comandos para iniciar tu experiencia:

• Envía <b>el enlace de stickers de LINE/Kakao</b> para importar o descargar 📦
• Envía <b>stickers, enlaces o GIFs de Telegram</b> para descargarlos o exportarlos a WhatsApp
• Envía <b>palabras clave</b> para buscar paquetes de stickers 🔎
• Usa <b>/create</b> o <b>/manage</b> para crear o administrar paquetes de stickers y CustomEmoji
• Usa <b>/command_list</b> para ver todos los comandos disponibles 👾

╔╗╔╗╔═╗╔╗─╔╗─╔═╗
║╚╝║║╦╝║║─║║─║║║
║╔╗║║╩╗║╚╗║╚╗║║║
╚╝╚╝╚═╝╚═╝╚═╝╚═╝`,
	}

	return c.Send(photo, tele.ModeHTML)
}

func sendCommandList(c tele.Context) error {
	message := `
• Comandos disponibles:

<b>/import</b>  <b>/search</b> → Importar o buscar stickers de LINE/Kakao ✅
<b>/download</b>  <b>/create</b>  <b>/manage</b> → Descargar, crear y administrar paquetes de stickers de Telegram ✅
`
	return c.Send(message, tele.ModeHTML, tele.NoPreview)
}

func sendAskEmoji(c tele.Context) error {
	selector := &tele.ReplyMarkup{}
	btnManu := selector.Data("Asignar por separado", "manual")
	btnRand := selector.Data(`Asignar en lote como "⭐"`, "random")
	selector.Inline(selector.Row(btnManu), selector.Row(btnRand))

	return c.Send(`Para las funciones, Telegram requiere que cada sticker posea un emoji y palabras clave para identificarlo:

• Pulsa "Asignar por separado" para asignar un emoji o palabras clave a cada uno.
• Envía un emoji para asignar el mismo a todos los stickers automáticamente.`,
		selector)
}

func sendConfirmExportToWA(c tele.Context, sn string, hex string) error {
	selector := &tele.ReplyMarkup{}
	baseUrl, _ := url.JoinPath(msbconf.WebappUrl, "export")
	webAppUrl := fmt.Sprintf("%s?sn=%s&hex=%s", baseUrl, sn, hex)
	log.Debugln("webapp export link is:", webAppUrl)
	webapp := tele.WebApp{URL: webAppUrl}
	btnExport := selector.WebApp("Continuar exportación →", &webapp)
	selector.Inline(selector.Row(btnExport))

	return c.Reply(`La exportación a WhatsApp requiere <a href="https://github.com/star-39/msb_app">Msb App</a> debido a sus restricciones. 
Luego pulsa "Continuar exportación".

Descarga:
<b>iPhone:</b> AppStore (N/A), <a href="https://github.com/star-39/msb_app/releases/latest/download/msb_app.ipa">IPA</a>
<b>Android:</b> GooglePlay (N/A), <a href="https://github.com/star-39/msb_app/releases/latest/download/msb_app.apk">APK</a>`,
		tele.ModeHTML, tele.NoPreview, selector)
}

func genSDnMnEInline(canManage bool, isTGS bool, _ string) *tele.ReplyMarkup {
	selector := &tele.ReplyMarkup{}
	btnSingle := selector.Data("Descargar este sticker", CB_DN_SINGLE)
	btnAll := selector.Data("Descargar paquete de stickers", CB_DN_WHOLE)
	btnMan := selector.Data("Gestionar paquete de stickers", CB_MANAGE)

	if canManage {
		selector.Inline(selector.Row(btnSingle), selector.Row(btnAll), selector.Row(btnMan))
	} else {
		if isTGS {
			// Si es TGS, no se soporta exportar a WA.
			selector.Inline(selector.Row(btnSingle), selector.Row(btnAll))
		} else {
			selector.Inline(selector.Row(btnSingle), selector.Row(btnAll))
		}
	}
	return selector
}

func sendAskSDownloadChoice(c tele.Context, s *tele.Sticker) error {
	selector := genSDnMnEInline(false, s.Animated, s.SetName)
	return c.Reply(`Puedes descargar este sticker o el paquete si lo necesitas. 
Por favor selecciona una opción:`,
		selector)
}

func sendAskSChoice(c tele.Context, sn string) error {
	selector := genSDnMnEInline(true, false, sn)
	return c.Reply(`Hey! Este paquete de stickers es tuyo 📚 
Puedes descargarlo o gestionarlo, selecciona una opción:`,
		selector)
}

func sendAskTGLinkChoice(c tele.Context) error {
	selector := &tele.ReplyMarkup{}
	btnManu := selector.Data("Descargar paquete de stickers", CB_DN_WHOLE)
	btnMan := selector.Data("Gestionar paquete de stickers", CB_MANAGE)
	selector.Inline(selector.Row(btnManu), selector.Row(btnMan))
	return c.Reply(`
Hey! Este paquete de stickers es tuyo 📚 
Puedes descargarlo o gestionarlo, selecciona una opción:`,
		selector)
}

func sendAskWantSDown(c tele.Context) error {
	selector := &tele.ReplyMarkup{}
	btn1 := selector.Data("Sí", CB_DN_WHOLE)
	btnNo := selector.Data("No", CB_BYE)
	selector.Inline(selector.Row(btn1), selector.Row(btnNo))
	return c.Reply(`
Puedes descargar este paquete de stickers. 
Pulsa "Sí" para continuar.`,
		selector)
}

func sendAskWantImportOrDownload(c tele.Context, avalAsEmoji bool) error {
	msg := ""
	selector := &tele.ReplyMarkup{}
	btnImportSticker := selector.Data("Importar como paquete de stickers", CB_OK_IMPORT)
	btnImportEmoji := selector.Data("Importar como Emoji personalizado", CB_OK_IMPORT_EMOJI)
	btnDownload := selector.Data("Descargar", CB_OK_DN)
	if avalAsEmoji {
		selector.Inline(selector.Row(btnImportSticker), selector.Row(btnImportEmoji), selector.Row(btnDownload))
		msg = `
Puedes importar este paquete de stickers a Telegram o descargarlo.
También puedes importarlo como Emoji personalizado, pero recuerda que necesitarás Telegram Premium para enviarlos.`
	} else {
		selector.Inline(selector.Row(btnImportSticker), selector.Row(btnDownload))
		msg = `
Puedes importar este paquete de stickers a Telegram o descargarlo.`
	}

	return c.Reply(msg, selector)
}

func sendAskWhatToDownload(c tele.Context) error {
	return c.Send("Por favor, envíame un sticker que quieras descargar o su link ^^ (puede ser de Telegram o LINE).\n")
}

func sendAskTitle_Import(c tele.Context) error {
	ld := users.data[c.Sender().ID].lineData
	ld.TitleWg.Wait()
	log.Debug("titles are::")
	log.Debugln(ld.I18nTitles)
	selector := &tele.ReplyMarkup{}

	var titleButtons []tele.Row
	var titleText string
	for i, t := range ld.I18nTitles {
		if t == "" {
			continue
		}
		title := escapeTagMark(t) + " @" + botName
		btn := selector.Data(title, strconv.Itoa(i))
		row := selector.Row(btn)
		titleButtons = append(titleButtons, row)
		titleText = titleText + "\n<code>" + title + "</code>"
	}

	if len(titleButtons) == 0 {
		btnDefault := selector.Data(escapeTagMark(ld.Title)+" @"+botName, CB_DEFAULT_TITLE)
		titleButtons = []tele.Row{selector.Row(btnDefault)}
	}
	selector.Inline(titleButtons...)

	return c.Send("Por favor envíame un título para este paquete de stickers 💫. También puedes seleccionar un título original:\n"+
		titleText, selector, tele.ModeHTML)
}

func sendAskTitle(c tele.Context) error {
	return c.Send("Por favor envíame un título para este paquete de stickers 💫.\n")
}

func sendAskID(c tele.Context) error {
	selector := &tele.ReplyMarkup{}
	btnAuto := selector.Data("Generar automáticamente 🔗", "auto")
	selector.Inline(selector.Row(btnAuto))
	return c.Send(`
Por favor envía un ID para el paquete de stickers, se usará para generar su link.
Solo puede contener letras, números y guiones bajos.

Por ejemplo:
<code>My_favSticker21</code>

El ID normalmente no es importante, puedes pulsar "Generar automáticamente".`, selector, tele.ModeHTML)
}

func sendAskImportLink(c tele.Context) error {
	return c.Send(`
Por favor envíame el enlace de la tienda LINE/Kakao del paquete de stickers ✅. 
Puedes obtener este enlace desde la aplicación, en la tienda de stickers, pulsando Compartir -> Copiar enlace.

Por ejemplo:
<code>https://store.line.me/stickershop/product/7673/ja</code>
<code>https://e.kakao.com/t/pretty-all-friends</code>
<code>https://emoticon.kakao.com/items/lV6K2fWmU7CpXlHcP9-ysQJx9rg=?referer=share_link</code>`,
		tele.ModeHTML)
}

func sendNotifySExist(c tele.Context, lineID string) bool {
	lines := queryLineS(lineID)
	if len(lines) == 0 {
		return false
	}
	message := "Hey! Este paquete de stickers ya existe en la base de datos. Puedes continuar con la importación o usarlo directamente si lo deseas.\n\n"

	var entries []string
	for _, l := range lines {
		if l.Ae {
			entries = append(entries, fmt.Sprintf(`<a href="%s">%s</a>`, "https://t.me/addstickers/"+l.Tg_id, l.Tg_title))
		} else {
			// añadir al inicio
			entries = append([]string{fmt.Sprintf(`★ <a href="%s">%s</a>`, "https://t.me/addstickers/"+l.Tg_id, l.Tg_title)}, entries...)
		}
	}
	if len(entries) > 5 {
		entries = entries[:5]
	}
	message += strings.Join(entries, "\n")
	c.Send(message, tele.ModeHTML)
	return true
}

func sendSearchResult(entriesWant int, lines []LineStickerQ, c tele.Context) error {
	var entries []string
	message := "Resultados de búsqueda 🔎:\n"

	for _, l := range lines {
		l.Tg_title = strings.TrimSuffix(l.Tg_title, " @"+botName)
		if l.Ae {
			entries = append(entries, fmt.Sprintf(`<a href="%s">%s</a>`, "https://t.me/addstickers/"+l.Tg_id, l.Tg_title))
		} else {
			// añadir al inicio
			entries = append([]string{fmt.Sprintf(`★ <a href="%s">%s</a>`, "https://t.me/addstickers/"+l.Tg_id, l.Tg_title)}, entries...)
		}
	}

	if entriesWant == -1 && len(entries) > 120 {
		c.Send("Demasiados resultados, por favor utiliza una palabra clave más precisa. Se han reducido a 120 entradas 😵.\n")
		entries = entries[:120]
	}
	if entriesWant != -1 && len(entries) > entriesWant {
		entries = entries[:entriesWant]
	}
	if len(entries) > 30 {
		eChunks := chunkSlice(entries, 30)
		for _, eChunk := range eChunks {
			msgToSend := message + strings.Join(eChunk, "\n")
			c.Send(msgToSend, tele.ModeHTML)
		}
	} else {
		message += strings.Join(entries, "\n")
		c.Send(message, tele.ModeHTML)
	}

	return nil
}

func sendAskStickerFile(c tele.Context) error {
	return c.Send("Empecemos!!! Por favor envíame imágenes/fotos/stickers (que sean menos de 120 en total) 📸,\n" +
		"o puedes enviarme un archivo comprimido que contenga las imágenes,\n" +
		"espera a que finalice la subida y luego pulsa 'Finalizar subida'.\n")
}

func sendInStateWarning(c tele.Context) error {
	command := users.data[c.Sender().ID].command
	state := users.data[c.Sender().ID].state

	return c.Send(fmt.Sprintf(`
Por favor envíame el contenido de acuerdo con las instrucciones ^^
Comando actual: %s
Estado actual: %s
También puedes usar /quit para cancelar la sesión.
`, command, state))
}

func sendNoSessionWarning(c tele.Context) error {
	return c.Send("Por favor usa /start o envía enlaces de LINE/Kakao/Telegram, o stickers\n")
}

func sendAskSTypeToCreate(c tele.Context) error {
	selector := &tele.ReplyMarkup{}
	btnRegular := selector.Data("Paquete de stickers ✨", CB_REGULAR_STICKER)
	btnCustomEmoji := selector.Data("Emojis personalizados 🌟", CB_CUSTOM_EMOJI)

	selector.Inline(selector.Row(btnRegular), selector.Row(btnCustomEmoji))
	return c.Send("¿Qué tipo de paquete quieres crear? ₍^. .^₎Ⳋ\nRecuerda que los emojis personalizados solo pueden ser enviados por usuarios con Telegram Premium.",
		selector)
}

func sendAskEmojiAssign(c tele.Context) error {
	sd := users.data[c.Sender().ID].stickerData
	sf := sd.stickers[sd.pos]
	sf.wg.Wait()
	caption := fmt.Sprintf(`
Envíame el o los emojis que representen este sticker ₍^. .^₎Ⳋ.

%d de %d
`, sd.pos+1, sd.lAmount)

	if sf.fileID != "" {
		msg, _ := c.Bot().Send(c.Sender(), &tele.Sticker{
			File: tele.File{FileID: sf.fileID},
		})
		_, err := c.Bot().Reply(msg, caption)
		return err
	}

	err := c.Send(&tele.Video{
		File:    tele.FromDisk(sf.oPath),
		Caption: caption,
	})
	if err != nil {
		err2 := c.Send(&tele.Video{
			File:    tele.FromDisk(sd.stickers[sd.pos].oPath),
			Caption: caption,
		})
		if err2 != nil {
			err3 := c.Send(&tele.Document{
				File:     tele.FromDisk(sd.stickers[sd.pos].oPath),
				FileName: filepath.Base(sd.stickers[sd.pos].oPath),
				Caption:  caption,
			})
			if err3 != nil {
				err4 := c.Send(&tele.Sticker{File: tele.File{FileID: sd.stickers[sd.pos].oPath}})
				if err4 != nil {
					return err4
				}
			}
		}
	}
	return nil
}

func sendFatalError(err error, c tele.Context) {
	if c == nil {
		return
	}
	var errMsg string
	if err != nil {
		errMsg = err.Error()
		errMsg = strings.ReplaceAll(errMsg, msbconf.BotToken, "***")
		if strings.Contains(errMsg, "500") {
			errMsg += "\nEste es un error interno del servidor de Telegram, no podemos hacer nada más que esperar a que se recupere. Por favor, inténtalo de nuevo más tarde ⏳\n"
		}
	}

	c.Send("<b>Error fatal. Por favor intenta de nuevo con /start</b>\n\n"+
		"Puedes reportar este error en https://github.com/star-39/moe-sticker-bot/issues\n\n"+
		"<code>"+errMsg+"</code>", tele.ModeHTML, tele.NoPreview)
}

func sendExecEmojiAssignFinished(c tele.Context) error {
	ud := users.data[c.Sender().ID]
	msg := fmt.Sprintf(`
🔸 LINE Cat: <code>%s</code>
🔸 LINE ID: <code>%s</code>
🔹 TG ID: <code>%s</code>
🔹 TG Title: <a href="%s">%s</a>

Completado con éxito ₍^. .^₎Ⳋ✅. /start
    `, ud.lineData.Category,
		ud.lineData.Id,
		ud.stickerData.id,
		"https://t.me/addstickers/"+ud.stickerData.id,
		escapeTagMark(ud.stickerData.title),
	)
	return c.Send(msg, tele.ModeHTML)
}

// Return:
// string: Text of the message.
// *tele.Message: The pointer of the message.
// error: error
func sendProcessStarted(ud *UserData, c tele.Context, optMsg string) (string, *tele.Message, error) {
	message := fmt.Sprintf(`
Preparando stickers, por favor espera ⏳...

🔸 LINE Cat: <code>%s</code>
🔸 LINE ID: <code>%s</code>
🔹 TG ID: <code>%s</code>
🔹 TG TYPE: <code>%s</code>
🔹 TG Title: <a href="%s">%s</a>

<b>Progreso en curso. . . ₍^. .^₎Ⳋ</b>
<code>%s</code>
`, ud.lineData.Category,
		ud.lineData.Id,
		ud.stickerData.id,
		ud.stickerData.stickerSetType,
		"https://t.me/addstickers/"+ud.stickerData.id,
		escapeTagMark(ud.stickerData.title),
		optMsg)
	ud.progress = message

	teleMsg, err := c.Bot().Send(c.Recipient(), message, tele.ModeHTML)
	ud.progressMsg = teleMsg
	return message, teleMsg, err
}

// if progressText is empty, a progress bar will be generated based on cur and total.
func editProgressMsg(cur int, total int, progressText string, originalText string, teleMsg *tele.Message, c tele.Context) error {
	defer func() {
		if r := recover(); r != nil {
			log.Errorln("editProgressMsg encountered panic! ignoring...", string(debug.Stack()))
		}
	}()

	header := originalText[:strings.LastIndex(originalText, "<code>")]
	prog := ""

	if progressText != "" {
		prog = progressText
		goto SEND
	}
	cur = cur + 1
	if cur == 1 {
		prog = fmt.Sprintf("<code>[=>                  ]\n       %d de %d</code>", cur, total)
	} else if cur == int(float64(0.25)*float64(total)) {
		prog = fmt.Sprintf("<code>[====>               ]\n       %d de %d</code>", cur, total)
	} else if cur == int(float64(0.5)*float64(total)) {
		prog = fmt.Sprintf("<code>[=========>          ]\n       %d de %d</code>", cur, total)
	} else if cur == int(float64(0.75)*float64(total)) {
		prog = fmt.Sprintf("<code>[==============>     ]\n       %d de %d</code>", cur, total)
	} else if cur == total {
		prog = fmt.Sprintf("<code>[====================]\n       %d de %d</code>", cur, total)
	} else {
		return nil
	}
SEND:
	messageText := header + prog
	c.Bot().Edit(teleMsg, messageText, tele.ModeHTML)
	return nil
}

func sendAskSToManage(c tele.Context) error {
	return c.Send("Envíame un sticker del paquete que quieras editar,\n" +
		"o envía su link para simplificar~ \n")
}

func sendUserOwnedS(c tele.Context) error {
	usq := queryUserS(c.Sender().ID)
	if usq == nil {
		return errors.New("Aun no tienes ningún paquete de stickers 🔴")
	}

	var entries []string

	for _, us := range usq {
		date := time.Unix(us.timestamp, 0).Format("2006-01-02 15:04")
		title := strings.TrimSuffix(us.tg_title, " @"+botName)
		// solución para título vacío
		if title == "" || title == " " {
			title = "_"
		}
		entry := fmt.Sprintf(`<a href="https://t.me/addstickers/%s">%s</a>`, us.tg_id, title)
		entry += " | " + date
		entries = append(entries, entry)
	}

	if len(entries) > 30 {
		eChunks := chunkSlice(entries, 30)
		for _, eChunk := range eChunks {
			message := "Tienes los siguientes paquetes de stickers:\n"
			message += strings.Join(eChunk, "\n")
			c.Send(message, tele.ModeHTML)
		}
	} else {
		message := "Estos son tus paquetes de stickers agregados:\n"
		message += strings.Join(entries, "\n")
		c.Send(message, tele.ModeHTML)
	}
	return nil
}

func sendAskEditChoice(c tele.Context) error {
	ud := users.data[c.Sender().ID]
	selector := &tele.ReplyMarkup{}
	btnAdd := selector.Data("Añadir sticker", CB_ADD_STICKER)
	btnDel := selector.Data("Eliminar sticker", CB_DELETE_STICKER)
	btnDelset := selector.Data("Eliminar paquete de stickers", CB_DELETE_STICKER_SET)
	btnChangeTitle := selector.Data("Cambiar título", CB_CHANGE_TITLE)
	btnExit := selector.Data("Salir", "bye bye~")

	if msbconf.WebappUrl != "" {
		baseUrl, _ := url.JoinPath(msbconf.WebappUrl, "edit")
		url := fmt.Sprintf("%s?ss=%s&dt=%d",
			baseUrl,
			ud.stickerData.id,
			time.Now().Unix())
		log.Debugln("WebApp URL es: ", url)
		webApp := &tele.WebApp{
			URL: url,
		}
		btnEdit := selector.WebApp("Cambiar orden o emoji", webApp)
		selector.Inline(
			selector.Row(btnAdd), selector.Row(btnDel), selector.Row(btnDelset), selector.Row(btnEdit), selector.Row(btnChangeTitle), selector.Row(btnExit))
	} else {
		selector.Inline(
			selector.Row(btnAdd), selector.Row(btnDel), selector.Row(btnDelset), selector.Row(btnChangeTitle), selector.Row(btnExit))
	}

	return c.Send(fmt.Sprintf(`
ID: <code>%s</code>
Título: <a href="https://t.me/addstickers/%s">%s</a>

¿Qué quieres editar? Selecciona una opción:`,
		users.data[c.Sender().ID].stickerData.id,
		ud.stickerData.id,
		ud.stickerData.title),
		selector, tele.ModeHTML)
}

func sendAskSDel(c tele.Context) error {
	return c.Send("¿Cuál sticker quisieras eliminar? Por favor envíalo.\n")
}

func sendConfirmDelset(c tele.Context) error {
	selector := &tele.ReplyMarkup{}
	btnYes := selector.Data("Sí", CB_YES)
	btnNo := selector.Data("No", CB_NO)
	selector.Inline(selector.Row(btnYes), selector.Row(btnNo))

	return c.Send("Estás intentando eliminar el paquete de stickers!! Por favor confirma.", selector)
}

func sendSFromSS(c tele.Context, ssid string, reply *tele.Message) error {
	ss, _ := c.Bot().StickerSet(ssid)
	if reply != nil {
		c.Bot().Reply(reply, &ss.Stickers[0])
	} else {
		c.Send(&ss.Stickers[0])
	}
	return nil
}

func sendFLWarning(c tele.Context) error {
	return c.Send(`
Puede tardar más en procesar este paquete de stickers (2-8 minutos)... 
Este aviso indica que probablemente activaste el límite de frecuencia de Telegram, y el bot está intentando reenviar.
Debido a este mecanismo, el paquete resultante puede contener stickers duplicados o faltantes, por favor revisa manualmente al finalizar.
`)
}

// func sendTooManyFloodLimits(c tele.Context) error {
//return c.Send("Lo siento, parece que activaste el límite de frecuencia de Telegram demasiadas veces. Se recomienda intentarlo de nuevo más tarde.\n")
//}

func sendNoCbWarn(c tele.Context) error {
	return c.Send("¡Por favor pulsa un botón! /quit\n")
}

func sendBadIDWarn(c tele.Context) error {
	return c.Send(`
ID inválido. Intenta de nuevo o pulsa Generar automáticamente o cancela usando /quit ₍^. .^₎Ⳋ
Solo puede contener letras, números y guiones bajos, debe comenzar con una letra y no puede tener guiones bajos consecutivos.

Por ejemplo:
<code>My_favSticker21</code>
`, tele.ModeHTML)
}

func sendIDOccupiedWarn(c tele.Context) error {
	return c.Send("El ID ya está ocupado. Intenta con otro ❌\n")
}

func sendBadImportLinkWarn(c tele.Context) error {
	return c.Send("Enlace de importación inválido, asegúrate de que sea un enlace de la tienda LINE o Kakao. Intenta de nuevo o /quit\n\n"+
		"Por ejemplo:\n"+
		"<code>https://store.line.me/stickershop/product/7673/ja</code>\n"+
		"<code>https://e.kakao.com/t/pretty-all-friends</code>", tele.ModeHTML)
}

func sendNoSToManage(c tele.Context) error {
	return c.Send("Lo siento, aún no has creado ningún paquete de stickers. Puedes usar /import o /create.\n")
}

func sendPromptStopAdding(c tele.Context) error {
	selector := &tele.ReplyMarkup{}
	btnDone := selector.Data("Finalizar subida", CB_DONE_ADDING)
	selector.Inline(selector.Row(btnDone))
	return c.Send("Continúa enviando archivos o pulsa el botón de abajo para finalizar.", selector)
}

func replySFileOK(c tele.Context, count int) error {
	selector := &tele.ReplyMarkup{}
	btnDone := selector.Data("Finalizar subida", CB_DONE_ADDING)
	selector.Inline(selector.Row(btnDone))
	return c.Reply(
		fmt.Sprintf("Archivo subido con éxito. Se recibieron %d stickers. Continúa enviando archivos o pulsa el botón de abajo para finalizar.", count), selector)
}

func sendSEditOK(c tele.Context) error {
	return c.Send("Paquete de stickers editado con éxito. /start")
}

func sendStickerSetFullWarning(c tele.Context) error {
	return c.Send("Aviso: Tu paquete de stickers ya está lleno. No puedes añadir más stickers.\n")
}

func sendAskSearchKeyword(c tele.Context) error {
	return c.Send("Por favor envía una palabra que quieras buscar.\n")
}

func sendSearchNoResult(c tele.Context) error {
	message := "Lo siento, no hay resultados."
	if c.Chat().Type == tele.ChatPrivate {
		message += "\nIntenta de nuevo o usa /quit"
	}
	return c.Send(message)
}

func sendNotifyNoSessionSearch(c tele.Context) error {
	return c.Send("Aquí tienes algunos resultados de búsqueda. Usa /search para profundizar más o /start para ver los comandos disponibles.\n")
}

func sendUnsupportedCommandForGroup(c tele.Context) error {
	return c.Send("Este comando no está soportado en chats de grupo, por favor habla directamente con el bot.\n")
}

func sendBadSearchKeyword(c tele.Context) error {
	return c.Send(fmt.Sprintf(`
Por favor especifica una palabra clave.

Ejemplo:
/search@%s palabra1 palabra2 ...
/search@%s nekomimi mia
`, botName, botName))
}

func sendPreferKakaoShareLinkWarning(c tele.Context) error {
	msg := `
El enlace que enviaste es de la tienda Kakao 🛍
Usa un enlace de compartir para mejorar la calidad de imagen y soportar stickers animados.
Puedes obtenerlo desde la app KakaoTalk en la tienda de stickers, pulsando Compartir -> Copiar enlace.

Ejemplo: <code>https://emoticon.kakao.com/items/lV6K2fWmU7CpXlHcP9-ysQJx9rg=?referer=share_link</code>
`
	err := c.Reply(&tele.Photo{
		File:    tele.File{FileID: FID_KAKAO_SHARE_LINK},
		Caption: msg,
	}, tele.ModeHTML)
	if err != nil {
		c.Reply(msg, tele.ModeHTML)
	}
	return nil
}

func sendUseCommandToImport(c tele.Context) error {
	return c.Send("Por favor usa /create para crear un paquete de stickers usando tus propias fotos y videos ₍^. .^₎Ⳋ /start\n")
}

func sendOneStickerFailedToAdd(c tele.Context, pos int, err error) error {
	return c.Reply(fmt.Sprintf(`
Error al añadir un sticker...
Índice: %d
Error: %s
`, pos, err.Error()))
}

func sendBadSNWarn(c tele.Context) error {
	return c.Reply("¡Sticker o enlace incorrecto! ❌\n")
}

func sendSSTitleChanged(c tele.Context) error {
	msg := `
Título cambiado con éxito✅ `
	return c.Reply(msg, tele.ModeHTML)
}

func sendSSTitleFailedToChanged(c tele.Context) error {
	msg := `
Error al cambiar el título, por favor inténtalo de nuevo.`
	return c.Reply(msg, tele.ModeHTML)
}

// func sendInvalidEmojiWarn(c tele.Context) error {
//  return c.Reply(`
// Lo siento, este emoji no es válido, se ha establecido por defecto como ⭐️.
// Puedes editarlo después usando el comando /manage.
//  `)
// }

func sendProcessingStickers(c tele.Context) error {
	return c.Send(`
Procesando stickers, por favor espera un momento...
`)
}
