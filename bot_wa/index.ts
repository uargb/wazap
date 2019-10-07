import { create, Whatsapp } from 'sulla-hotfix'
import axios from 'axios'
import * as fs from 'fs'
import * as mime from 'mime-types'

function base64_encode(file: string): string {
  var data = fs.readFileSync(file);
  return Buffer.from(data).toString('base64');
}

function apiBase(phone: string, url: string) { 
  return 'http://127.0.0.1:8090/bot/' + phone + '/' + url
}


async function start(client: Whatsapp) {
  client.onMessage(async message => {
    const phone = message.chatId.split('@')[0]
    try {
      let response = await axios.get(apiBase(phone, 'answer?message=' + encodeURIComponent(message.body)))
      if (response.data.ok) {
        if (response.data.did === 'registered') {
          try {
            await axios.get(apiBase(phone, 'rename?name=' + encodeURIComponent(message.sender.pushname)))
          } catch (error) {
            console.error(error)
          }
        }

        response.data.data.forEach(async (card: any) => {
          let mediaData = 'data:{mime};base64,{base64}'
          let mediaName = ''

          if (card.Image) {
            mediaName = card.Image
            mediaData = mediaData.replace('{mime}', mime.lookup(card.Image))
            mediaData = mediaData.replace('{base64}', base64_encode('./public/' + card.ManagerID + '-' + card.Image))
          } else if (card.Video) {
            mediaName = card.Video
            mediaData = mediaData.replace('{mime}', 'application/octet-stream')
            mediaData = mediaData.replace('{base64}', base64_encode('./public/' + card.ManagerID + '-' + card.Video))
          } else if (card.Attachment) {
            mediaName = card.Attachment
            mediaData = mediaData.replace('{mime}', mime.lookup(card.Attachment))
            mediaData = mediaData.replace('{base64}', base64_encode('./public/' + card.ManagerID + '-' + card.Attachment))
          }

          if (mediaName.length > 0) {
            await client.sendImage(message.from, mediaData, mediaName, card.Text)
          } else {
            await client.sendText(message.from, card.Text)
          }

          if (card.NotifyManager) {
            await client.sendText(response.data.manager.Phone + '@c.us', JSON.stringify(response.data.costumer))
          }
        })
      }
    } catch (error) {
      console.error(error)
    }
  })
}


create().then((client: Whatsapp) => {
  start(client)
}).catch((error: Error) => console.error(error))