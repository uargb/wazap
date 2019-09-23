<template>
  <section id="menu">
    <navbar :username="username" :password="password" />
    <br />
    <div id="menu" class="container">
      <div class="columns is-multiline">
        <div class="column is-4" v-for="(item, index) in data" :key="index">
          <div class="box is-light">
            <div class="buttons justify-space-between has-margin-bottom-20">
                <b-tooltip label="Добавить новый пункт">
                  <b-button type="is-success"
                    icon-left="arrow-left" @click="prepend(index)"></b-button>
                  <b-button type="is-success"
                    icon-left="arrow-right" @click="append(index)"></b-button>
                </b-tooltip>

              <b-tooltip label="Удалить пункт">
                <b-button type="is-danger"
                icon-left="delete"
                @click="remove(index)"
                :loading="item.loading"
                :disabled="item.loading"></b-button>
              </b-tooltip>
            </div>
            <b-field label="Запрос" label-position="inside">
              <b-input v-model="item.query"></b-input>
            </b-field>
            <b-field label="Описание запроса" label-position="inside">
              <b-input v-model="item.description"></b-input>
            </b-field>
            <div class="field">
              <b-checkbox v-model="item.show">Отображать в приветствии</b-checkbox>
            </div>
            <b-field label="Текст" label-position="inside">
              <b-input type="textarea" v-model="item.text"></b-input>
            </b-field>
            <b-field class="file">
              <b-upload v-model="item.image">
                  <a class="button is-secondary">
                      <b-icon icon="upload"></b-icon>
                      <span>Картинка</span>
                  </a>
              </b-upload>
              <span class="file-name" v-if="item.image.name">
                  {{ item.image.name }}
              </span>
            </b-field>
            <b-field class="file">
              <b-upload v-model="item.video">
                  <a class="button is-secondary">
                      <b-icon icon="upload"></b-icon>
                      <span>Видео</span>
                  </a>
              </b-upload>
              <span class="file-name" v-if="item.video.name">
                  {{ item.video.name }}
              </span>
            </b-field>
            <b-field class="file">
              <b-upload v-model="item.attachment">
                  <a class="button is-secondary">
                      <b-icon icon="upload"></b-icon>
                      <span>Вложение</span>
                  </a>
              </b-upload>
              <span class="file-name" v-if="item.attachment.name">
                  {{ item.attachment.name }}
              </span>
            </b-field>
            <b-button
              type="is-info"
              @click="save(index)"
              :loading="item.loading"
              :disabled="item.loading"
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
    username () {
      return this.$store.state.username
    },
    password () {
      return this.$store.state.password
    }
  },
  mounted () {
    this.load()
  },
  methods: {
    async load () {
      this.pageLoading = true
      try {
        let response = await this.$axios.get('qa', { params: {
          username: this.username,
          password: this.password
        }})
        if (response.data.ok) {
          let payload = []
          if (response.data.data > 2) {
            payload = JSON.parse(response.data.data)
          }
          if (payload.length == 0) {
            this.data.push({
              loading: false,
              query: '',
              description: '',
              show: true,
              text: '',
              image: { name: '' },
              video: { name: '' },
              attachment: { name: '' }
            })
          } else {
            payload.forEach(item => {
              this.data.push({
                loading: false,
                query: item.query,
                description: item.description,
                show: payload.show > 0,
                text: item.text,
                image: { name: item.image },
                video: { name: item.video },
                attachment: { name: item.attachment }
              })
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
    async prepend(index) {
      this.data.splice(index, 0, {
        loading: false,
        query: '',
        description: '',
        show: true,
        text: '',
        image: {},
        video: {},
        attachment: {}
      })

      this.pageLoading = true
      for (let i = index; i < this.data.length; i++) {
        await this.save(i)
      }
      this.pageLoading = false
    },
    async append (index) {
      this.data.splice(index + 1, 0, {
        loading: false,
        query: '',
        description: '',
        show: true,
        text: '',
        image: {},
        video: {},
        attachment: {}
      })
      await this.save(index + 1)
    },
    async remove(index) {
      this.data[index].loading = true
      try {
        let response = await this.$axios.post('/qa/remove', this.$qs.stringify({
          username: this.username,
          password: this.password,
          index: index
        }))
        if (response.data.ok) {
          this.data.splice(index, 1)
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
        this.data[index].loading = false
      }
    },
    async save (index) {
      this.data[index].loading = true
      try {
        let data = new FormData();
        data.append('username', this.username)
        data.append('password', this.password)
        data.append('index', index)
        data.append('query', this.data[index].query)
        data.append('description', this.data[index].description)
        data.append('show', this.data[index].show ? 1 : 0)
        data.append('text', this.data[index].text)
        if (this.data[index].image instanceof File) {
          data.append('image', this.data[index].image)
        }
        if (this.data[index].video instanceof File) {
          data.append('video', this.data[index].video)
        }
        if (this.data[index].attachment instanceof File) {
          data.append('attachment', this.data[index].attachment)
        }

        let response = await this.$axios.post('/qa', data, {
          headers: {"Content-Type": "multipart/form-data"},
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
        this.data[index].loading = false
      }
    }
  }
}
</script>
