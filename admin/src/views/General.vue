<template>
  <section id="general">
    <navbar />
    <br/>
    <div class="container">
      <div class="box">
        <b-field label="Имя">
          <b-input type="text" v-model="info.Name"></b-input>
        </b-field>
        <b-field label="Шаблон текста рекламной ссылки">
          <b-input type="text" v-model="info.LinkTemplate"></b-input>
        </b-field>
        <b-field label="Рекламная ссылка">
          <b-input type="text" v-model="info.Link" readonly></b-input>
        </b-field>
        <b-field label="Приветствие">
          <b-input type="textarea" v-model="info.Greeting"></b-input>
        </b-field>
        <b-field label="Номер телефона">
          <b-input v-model="info.Phone"></b-input>
        </b-field>
        <b-button type="is-info"
          @click="save"
          :loading="loading" :disabled="loading"
          >Сохранить</b-button>
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
  name: 'general',
  components: { Navbar },
  data () {
    return {
      loading: false,
      info: {}
    }
  },
  computed: {
    username () { return this.$store.state.username },
    password () { return this.$store.state.password }
  },
  mounted () {
    this.load()
  },
  methods: {
    async load () {
      this.loading = true
      try {
        let response = await axios.get(this.$apiBase(this.username, this.password, 'general'))
        if (response.data.ok) {
          this.info = response.data.data
        } else {
          this.$error(this, response.data.message)
        }
      } catch (error) {
        this.$error(this, error.message)
      } finally {
        this.loading = false
      }
    },
    async save () {
      this.loading = true
      try {
        let response = await axios.post(
          this.$apiBase(this.username, this.password, 'general'),
          this.$qs.stringify(this.info)
        )
        if (!response.data.ok) {
          this.$error(this, response.data.message)
        }
      } catch (error) {
        this.$error(this, error.message)
      } finally {
        this.load()
      }
    }
  }
}
</script>
