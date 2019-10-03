<template>
  <section id="login-page" class="hero is-primary is-fullheight">
    <div class="hero-body">
      <div class="container">
        <div class="column is-4 is-offset-4">
          <div class="box">
            <b-field label="Имя пользователя">
              <b-input type="text" v-model="username"></b-input>
            </b-field>
            <b-field label="Пароль">
              <b-input type="password" v-model="password"></b-input>
            </b-field>
            <b-button type="is-primary"
              @click="authorize"
              :loading="loading" :disabled="loading"
            >Вход</b-button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import axios from 'axios'
export default {
  name: 'login',
  data () {
    return {
      loading: false,
      username: '',
      password: ''
    }
  },
  created () { this.$store.commit('logout') },
  methods: {
    async authorize () {
      this.loading = true
      try {
        let response = await axios.post(this.$apiBase(this.username, this.password, 'check'))
        if (response.data.ok) {
          this.$store.commit('authenticate', {
            username: this.username,
            password: this.password
          })
          this.$router.push('general')
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
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
