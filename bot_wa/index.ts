import { create, Whatsapp } from 'sulla-hotfix'
import * as mysql from 'promise-mysql'
import * as fs from 'fs'
import * as mime from 'mime-types'

const dbConfig = {
  host: 'wazap.cvtuyrclurh0.ap-south-1.rds.amazonaws.com',
  post: 3306,
  user: 'admin',
  password: '952368741',
  database: 'gym'
}


function translateQuery(query: string): string {
  return query
}

function base64_encode(file: string): string {
  var data = fs.readFileSync(file);
  return Buffer.from(data).toString('base64');
}


async function start(db: mysql.Connection, client: Whatsapp) {
  client.onMessage(async message => {
    const phone = message.chatId.split('@')[0]
    let result = await db.query('select * from costumers where phone = ?', [phone])
    if (result.length == 0) {
      // зарегистирурем пользователя как нового
      let managers = await db.query('select id, name, greeting from managers')
      
      // по умолчанию скидываем админу
      let manager = (await db.query('select id, greeting, qa from managers where id = 1'))[0]  
      managers.forEach(async (m: any) => {
        if (message.body.includes(manager.name)) {
          manager = m // если в сообщении есть имя менеджера - кинем ему
        }
      })

      await db.query(
        'insert into costumers (phone, name, managerId) values(?, ?, ?)',
        [phone, message.sender.pushname, manager.id]
      )

      let response = manager.greeting + "\n\n"
      const qa = JSON.parse(manager.qa)
      qa.forEach((item: any) => {
        if (item.show > 0) {
          response += translateQuery(item.query) + ' - ' + item.description + "\n"
        }
      });

      await client.sendText(message.from, response)
      return
    }
    
    const user = result[0]
    const manager = (await db.query(
      'select * from managers where id = ?',
      user.managerId)
    )[0]
    const qa = JSON.parse(manager.qa)

    let selected = [] // все карточки по выбранному запросу
    let unknown: any = undefined // карточка "неопознанного" запроса
    qa.forEach((item: any) => {
      if (item.query === '<неизвестный>')
        unknown = item

      if (item.query.includes(message.body.trim()))
        selected.push(item)
    })

    if (selected.length == 0 && unknown) 
      selected.push(unknown)

    selected.forEach(async (item) => {
      let mediaData = 'data:{mime};base64,{base64}'
      let mediaName = ''
      if (item.image) {
        mediaName = item.image
        mediaData = mediaData.replace('{mime}', mime.lookup(item.image))
        mediaData = mediaData.replace('{base64}', base64_encode('./public/' + item.image))
      } else if (item.video) {
        mediaName = item.video
        mediaData = mediaData.replace('{mime}', 'application/octet-stream')
        mediaData = mediaData.replace('{base64}', base64_encode('./public/' + item.video))
      } else if (item.attachment) {
        mediaName = item.attachment
        mediaData = mediaData.replace('{mime}', mime.lookup(item.attachment))
        mediaData = mediaData.replace('{base64}', base64_encode('./public/' + item.attachment))
      }
      if (mediaName.length > 0)
        client.sendImage(message.from, mediaData, mediaName, item.text)
      else
        client.sendText(message.from, item.text)
    })
  })
}


mysql.createConnection(dbConfig).then((db: mysql.Connection) => {
  create().then((client: Whatsapp) => {
    start(db, client)
  }).catch((error: Error) => { throw error })
}).catch((error: Error) => console.error(error))