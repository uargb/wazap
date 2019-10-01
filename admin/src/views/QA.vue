<template>
  <section id="menu">
    <navbar :username="username" :password="password" />
    <br />
    <div id="menu" class="container">
      <div class="columns is-multiline">
        <div class="column is-4" v-for="card in cards" :key="card.id">
          <div class="box is-light">
            <div class="buttons justify-space-between has-margin-bottom-20">
                <b-tooltip label="Добавить новый пункт">
                  <b-button type="is-success"
                    icon-left="arrow-left" @click="prepend(card.id)"></b-button>
                  <b-button type="is-success"
                    icon-left="arrow-right" @click="append(card.id)"></b-button>
                </b-tooltip>

              <b-tooltip label="Удалить пункт">
                <b-button type="is-danger"
                icon-left="delete"
                @click="remove(card.id)"
                :loading="card.loading"
                :disabled="card.loading"></b-button>
              </b-tooltip>
            </div>
            <b-field label="Запрос" label-position="inside">
              <b-input v-model="card.query"></b-input>
            </b-field>
            <b-field label="Описание запроса" label-position="inside">
              <b-input v-model="card.description"></b-input>
            </b-field>
            <div class="field">
              <b-checkbox v-model="card.show">Отображать в приветствии</b-checkbox>
            </div>
            <b-field label="Текст" label-position="inside">
              <b-input type="textarea" v-model="card.text"></b-input>
            </b-field>
            <b-field class="file">
              <b-upload v-model="card.image">
                  <a class="button is-secondary">
                      <b-icon icon="upload"></b-icon>
                      <span>Картинка</span>
                  </a>
              </b-upload>
              <span class="file-name" v-if="card.image.name">
                  {{ card.image.name }}
              </span>
            </b-field>
            <b-field class="file">
              <b-upload v-model="card.video">
                  <a class="button is-secondary">
                      <b-icon icon="upload"></b-icon>
                      <span>Видео</span>
                  </a>
              </b-upload>
              <span class="file-name" v-if="card.video.name">
                  {{ card.video.name }}
              </span>
            </b-field>
            <b-field class="file">
              <b-upload v-model="card.attachment">
                  <a class="button is-secondary">
                      <b-icon icon="upload"></b-icon>
                      <span>Вложение</span>
                  </a>
              </b-upload>
              <span class="file-name" v-if="card.attachment.name">
                  {{ card.attachment.name }}
              </span>
            </b-field>
            <b-button
              type="is-info"
              @click="save(card.id)"
              :loading="card.loading"
              :disabled="card.loading"
            >Сохранить</b-button>
          </div>
        </div>
      </div>
    </div>
    <br/>
    <b-loading is-full-page :active="pageLoading"></b-loading>
  </section>
</template>

<script>
import Navbar from '@/components/Navbar'
export default {
  name: 'qa',
  components: { Navbar },
  data () {
    return {
      pageLoading: false,
      data: []
    }
  },
  computed: {
    username () { return this.$store.state.username },
    password () { return this.$store.state.password }
  },
  mounted () { this.load() },
  methods: {
    newCard(data) {
      return {
        loading: false,
        query: ('query' in data) ? data['query'] : '',
        description: ('description' in data) ? data['description'] : '',
        show: true,
        text: ('text' in data) ? data['text'] : '',
        image: { name: ('image' in data) ? data['image'] : '' },
        video: { name: ('video' in data) ? data['video'] : '' },
        attachment: { name: ('attachment' in data) ? data['attachment'] : '' }
      }
    },
    async load () {
      this.pageLoading = true
      try {
        let response = await this.$axios.get('qa', { params: {
          username: this.username,
          password: this.password
        } })
        if (response.data.ok) {
          let payload = []
          if (response.data.data.length > 2) {
            payload = JSON.parse(response.data.data)
          }
          if (payload.length === 0) {
            this.data.push(this.newCard())
          } else {
            payload.forEach(card => {
              this.data.push(this.newCard(card))
            })
          }
        } else {
          this.$buefy.toast.open({
            message: response.data.message,
            type: 'is-danger',
            duration: 2000
          })
          this.$router.push('/')
        }
      } catch (error) {
        this.$buefy.toast.open({
          message: error.message,
          type: 'is-warning',
          duration: 2000
        })
      } finally {
        this.pageLoading = false
      }
    },
    async prepend (id) {
      this.data.splice(id, 0, this.newCard())

      this.pageLoading = true
      for (let i = id; i < this.data.length; i++) {
        await this.save(i)
      }
      this.pageLoading = false
    },
    async append (id) {
      this.data.splice(id + 1, 0, this.newCard())
      await this.save(id + 1)
    },
    async remove (id) {
      this.data[id].loading = true
      try {
        let response = await this.$axios.post('/qa/remove', this.$qs.stringify({
          username: this.username,
          password: this.password,
          id: id
        }))
        if (response.data.ok) {
          this.data.splice(id, 1)
        } else {
          this.$buefy.toast.open({
            message: response.data.message,
            type: 'is-danger',
            duration: 2000
          })
        }
      } catch (error) {
        this.$buefy.toast.open({
          message: error.message,
          type: 'is-warning',
          duration: 2000
        })
        this.data[id].loading = false
      }
    },
    async save (id) {
      this.data[id].loading = true
      try {
        let data = new FormData()
        data.append('username', this.username)
        data.append('password', this.password)
        for (var key in this.data) {
          data.append(key, this.data[id][key]);
        }

        let response = await this.$axios.post('/qa', data, {
          headers: { 'Content-Type': 'multipart/form-data' }
        })
        if (!response.data.ok) {
          this.$buefy.toast.open({
            message: response.data.message,
            type: 'is-danger',
            duration: 2000
          })
        }
      } catch (error) {
        this.$buefy.toast.open({
          message: error.message,
          type: 'is-warning',
          duration: 2000
        })
      } finally {
        this.data[id].loading = false
      }
    }
  }
}
</script>
