<template>
  <section id="menu">
    <navbar :username="username" :password="password" />
    <br />
    <div id="menu" class="container">
      <div class="columns">
        <div class="column is-12">
          <div class="box flex-row align-items-center">
            <p class="heading has-margin-right-15">Управление</p>
            <b-button type="is-success" icon-left="plus" @click="add()">Добавить новую карточку</b-button>
          </div>
        </div>
      </div>
      <div class="columns is-multiline">
        <div class="column is-4" v-for="card in cards" :key="card.ID">
          <div class="box is-light" :class="{'has-background-dark': card.disabled}">
            <div class="buttons justify-space-between has-margin-bottom-20">
              <div></div>
              <b-tooltip label="Удалить пункт">
                <b-button type="is-danger"
                icon-left="delete"
                @click="remove(card.ID)"
                :loading="card.loading"
                :disabled="card.loading || card.disabled"></b-button>
              </b-tooltip>
            </div>
            <b-field label="ID" label-position="inside">
              <b-input disabled v-model="card.ID"></b-input>
            </b-field>
            <b-field label="Запрос" label-position="inside">
              <b-input v-model="card.Query" :disabled="card.disabled"></b-input>
            </b-field>
            <b-field label="Текст" label-position="inside">
              <b-input type="textarea" v-model="card.Text" :disabled="card.disabled"></b-input>
            </b-field>
            <b-field label="Картинка" label-position="inside">
              <b-select expanded v-model="card.Image" :disabled="card.disabled">
                <option value="">Отсутствует</option>
                <option v-for="(image, index) in images" :key="index"
                  :value="image">{{ image }}</option>
              </b-select>
            </b-field>
            <b-field label="Видео" label-position="inside">
              <b-select expanded v-model="card.Video" :disabled="card.disabled">
                <option value="">Отсутствует</option>
                <option v-for="(video, index) in videos" :key="index" :value="video">{{ video }}</option>
              </b-select>
            </b-field>
            <b-field label="Вложение" label-position="inside">
              <b-select expanded v-model="card.Attachment" :disabled="card.disabled">
                <option value="">Отсутствует</option>
                <option v-for="(attachment, index) in attachments" :key="index" :value="attachment">{{ attachment }}</option>
              </b-select>
            </b-field>
            <b-field label="Сохранить информацию" label-position="inside">
              <b-input v-model="card.Write" :disabled="card.disabled"></b-input>
            </b-field>
            <b-field label="Отправить независимо от ответа" label-position="inside">
              <b-select expanded v-model="card.Attachment" :disabled="card.disabled">
                <option value="">Отсутствует</option>
                <option v-for="(card, index) in cards" :key="index" :value="card.ID">#{{ card.ID }}</option>
              </b-select>
            </b-field>
            <b-button
              type="is-info"
              @click="save(card.ID)"
              :loading="card.loading"
              :disabled="card.loading || card.disabled"
            >Сохранить</b-button>
          </div>
        </div>
      </div>
    </div>
    <br/>
    <b-loading is-full-page :active="loading"></b-loading>
  </section>
</template>

<script>
import axios from 'axios'
import Navbar from '@/components/Navbar'
export default {
  name: 'qa',
  components: { Navbar },
  data () {
    return {
      loading: false,
      cards: [],
      images: [],
      videos: [],
      attachments: []
    }
  },
  computed: {
    username () { return this.$store.state.username },
    password () { return this.$store.state.password }
  },
  mounted () { this.load() },
  methods: {
    async load () {
      this.loading = true
      try {
        let response = await axios.get(this.$apiBase(this.username, this.password, 'qa'))
        if (response.data.ok) {
          this.cards = response.data.data
        } else {
          this.$error(this, response.data.message)
        }

        response = await axios.get(this.$apiBase(this.username, this.password, 'files?ext=jpg'))
        if (response.data.ok) {
          this.images = response.data.data
        } else {
          this.$error(this, response.data.message)
        }

        response = await axios.get(this.$apiBase(this.username, this.password, 'files?ext=mp4'))
        if (response.data.ok) {
          this.videos = response.data.data
        } else {
          this.$error(this, response.data.message)
        }

        response = await axios.get(this.$apiBase(this.username, this.password, 'files?ext=pdf'))
        if (response.data.ok) {
          this.attachments = response.data.data
        } else {
          this.$error(this, response.data.message)
        }
      } catch (error) {
        this.$error(this, error.message)
      } finally {
        this.loading = false
      }
    },
    async add () {
      this.loading = true
      try {
        let response = await axios.post(this.$apiBase(this.username, this.password, 'qa/create'))
        if (response.data.ok) {
          this.cards.push(response.data.data)
        } else {
          this.$error(this, response.data.message)
        }
      } catch (error) {
        this.$error(this, error.message)
      } finally {
        this.loading = false
      }
    },
    async remove (id) {
      this.loading = true
      try {
        let response = await axios.post(this.$apiBase(this.username, this.password, 'qa/remove?id=' + id))
        if (response.data.ok) {
          for (let i = 0; i < this.cards.length; i++) {
            if (this.cards[i].ID === id) {
              this.cards[i].disabled = true
              break
            }
          }
        } else {
          this.$error(this, response.data.message)
        }
      } catch (error) {
        this.$error(this, error.message)
      } finally {
        this.loading = false
      }
    },
    async save (id) {
      this.loading = true
      try {
        let response
        for (let i = 0; i < this.cards.length; i++) {
          if (this.cards[i].ID === id) {
            response = await axios.post(
              this.$apiBase(this.username, this.password, 'qa/modify?id=' + id),
              this.$qs.stringify(this.cards[i])
            )
            break
          }
        }

        if (response.data.ok) {
          response = await axios.get(this.$apiBase(this.username, this.password, 'qa?id=' + id))
          if (response.data.ok) {
            for (let i = 0; i < this.cards.length; i++) {
              if (this.cards[i].ID === id) {
                this.cards[i] = response.data.data
              }
            }
          } else {
            this.$error(this, response.data.message)
          }
        } else {
          this.$error(this, response.data.message)
        }
      } catch (error) {
        this.$error(this, error.message)
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
