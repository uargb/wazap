<template>
  <section id="general">
    <navbar :username="username" :password="password"/>
    <br/>
    <div class="container">
      <div id="basic" class="box">
        <b-field label="Имя">
          <b-input type="text" v-model="name"></b-input>
        </b-field>
        <b-field label="Шаблон текста рекламной ссылки">
          <b-input type="text" v-model="linkTemplate"></b-input>
        </b-field>
        <b-field label="Приветствие">
          <b-input type="textarea" v-model="greeting"></b-input>
        </b-field>
        <b-field label="Рекламная ссылка">
          <b-input type="text" v-model="link" readonly></b-input>
        </b-field>
        <b-button type="is-info"
          @click="save"
          :loading="loading" :disabled="loading"
          >Сохранить</b-button>
      </div>
    </div>
    <br/>
    <b-loading is-full-page :active="pageLoading"></b-loading>
  </section>
</template>

<script>
import Navbar from '@/components/Navbar'
export default {
  name: 'general',
  components: { Navbar },
  data () {
    return {
      pageLoading: false,
      loading: false,
      name: '',
      linkTemplate: '',
      greeting: '',
      link: ''
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
        let response = await this.$axios.get('general', { params: {
          username: this.username,
          password: this.password
        } })
        if (response.data.ok) {
          let payload = JSON.parse(response.data.data)
          this.name = payload.name
          this.linkTemplate = payload.linkTemplate
          this.greeting = payload.greeting
          this.link = payload.link
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
        this.$router.push('/')
      } finally {
        this.pageLoading = false
      }
    },
    async save () {
      this.loading = true
      try {
        let response = await this.$axios.post('general', this.$qs.stringify({
          username: this.username,
          password: this.password,
          name: this.name,
          linkTemplate: this.linkTemplate,
          greeting: this.greeting
        }))
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
        this.loading = false
      }
    }
  }
}
</script>
