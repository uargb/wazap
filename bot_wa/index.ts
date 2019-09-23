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


function base64_encode(file: string): string {
  var data = fs.readFileSync(file);
  return Buffer.from(data).toString('base64');
}

function start(db: mysql.Connection, client: Whatsapp) {
  client.onMessage(async message => {
    const phone = message.chatId.split('@')[0]
    let result = await db.query('select * from costumers where phone = ?', [phone])
    if (result.length == 0) {
      // зарегистирурем пользователя как нового
      let managers = await db.query('select id, name, greeting from managers')
      // по умолчанию скидываем админу
      let newManager = (await db.query('select id, greeting from managers where id = 1'))[0]  
      managers.forEach(async (manager: any) => {
        if (message.body.includes(manager.name)) {
          newManager = manager // если в сообщении есть имя менеджера - кинем ему
        }
      })

      await db.query(
        'insert into costumers (phone, name, managerId) values(?, ?, ?)',
        [phone, message.sender.pushname, newManager.id]
      )

      let response = newManager.greeting
      // todo: организовать меню

      client.sendText(message.from, response)
      return
    }
    
    const user = result[0]
    const manager = (await db.query(
      'select * from managers where id = ?',
      user.managerId)
    )[0]

    const qa = JSON.parse(manager.qa)
    let selected: any
    qa.some((item: any) => {
      if (item.query === '<неизвестный>') {
        selected = item
      }

      if (item.query.includes(message.body.trim())) {
        selected = item
        return true
      }
    })

    if (selected.image) {
      const data = 'data:' + mime.lookup(selected.image) + ';base64,' + base64_encode('./public/' + selected.image)
      await client.sendImage(
        message.from,
        data,
        selected.image,
        selected.text
      )
    } else if (selected.video) {

    } else if (selected.attachment) {
      const data = 'data:' + mime.lookup(selected.attachment) + ';base64,' + base64_encode('./public/' + selected.attachment)
      await client.sendImage(
        message.from,
        data,
        selected.attachment,
        selected.text
      )
    } else {
        client.sendText(message.from, selected.text)
    }
  })
}


mysql.createConnection(dbConfig).then((db: mysql.Connection) => {
  create().then((client: Whatsapp) => {
    start(db, client)
  }).catch((error: Error) => { throw error })
}).catch((error: Error) => console.error(error))