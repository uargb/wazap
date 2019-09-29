<template>
  <section id="general">
    <navbar :username="username" :password="password"/>
    <br/>
    <div class="container">
      <div class="columns is-multiline">
        <div class="column is-4" v-for="(item, index) in data" :key="index">
          <div class="box">
            <b-field label="Имя" label-position="inside">
              <b-input v-model="item.name" readonly></b-input>
            </b-field>
            <b-field label="Телефон" label-position="inside">
              <b-input v-model="item.phone" readonly></b-input>
            </b-field>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import Navbar from '@/components/Navbar'
export default {
  name: 'costumers',
  components: { Navbar },
  data () {
    return {
      pageLoging: false,
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
        let response = await this.$axios.get('costumers', { params: {
          username: this.username,
          password: this.password
        } })
        if (response.data.ok) {
          let payload = JSON.parse(response.data.data)
          payload.forEach(item => {
            this.data.push({
              name: item.name,
              phone: item.phone
            })
          })
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
    }
  }
}
</script>
