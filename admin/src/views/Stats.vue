<template>
  <section id="stats">
    <navbar />
    <br/>
    <div class="container">
      <div class="columns">
        <div class="column is-12">
          <div class="box flex-row align-items-center">
            <p class="heading has-margin-right-15 has-margin-top-5">Выберите менеджера</p>
            <b-select expanded v-model="selectedManager" @input="load()">
              <option value="-1">Отсутствует</option>
              <option v-for="manager in managers" :key="manager.ID" :value="manager.ID">{{ manager.Name }}</option>
            </b-select>
          </div>
        </div>
      </div>
      <div class="columns" v-show="selectedManager > 0">
        <div class="column is-6">
          <div class="box">
            <apexchart type=donut :series="dataset" :options="{labels: labels}"/>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import axios from 'axios'
import Navbar from '@/components/Navbar'
export default {
  name: 'stats',
  components: { Navbar },
  data () {
    return {
      loading: false,
      managers: [],
      selectedManager: -1,
      labels: [],
      dataset: []
    }
  },
  computed: {
    username () { return this.$store.state.username },
    password () { return this.$store.state.password }
  },
  async mounted () {
    this.loading = true
    try {
      let response = await axios.get(this.$apiBase(this.username, this.password, 'managers'))
      if (response.data.ok) {
        this.managers = response.data.data
      } else {
        this.$error(this, response.data.message)
      }
    } catch (error) {
      this.$error(this, error.message)
    } finally {
      this.loading = false
    }
  },
  methods: {
    async load () {
      this.loading = true
      try {
        let response = await axios.get(this.$apiBase(this.username, this.password, 'stats?id=' + this.selectedManager))
        if (response.data.ok) {
          this.labels = []
          this.dataset = []

          for (var key in response.data.data.stats) {
            if (key === '') {
              this.labels.push('Без статуса')
            } else {
              this.labels.push(key)
            }

            this.dataset.push(response.data.data.stats[key])
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
