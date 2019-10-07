<template>
  <section id="costumers">
    <navbar :username="username" :password="password"/>
    <br/>
    <div class="container">
      <div class="columns is-multiline">
        <div class="column is-3" v-for="costumer in costumers" :key="costumer.ID">
          <div class="box">
            <b-field label="ID" label-position="inside">
              <b-input v-model="costumer.ID" disabled></b-input>
            </b-field>
            <b-field label="Имя" label-position="inside">
              <b-input v-model="costumer.Name" readonly></b-input>
            </b-field>
            <b-field label="Телефон" label-position="inside">
              <b-input v-model="costumer.Phone" readonly></b-input>
            </b-field>
            <b-field label="Статус" label-position="inside">
              <b-input v-model="costumer.Status" readonly></b-input>
            </b-field>
            <b-field label="Данные" label-position="inside">
              <b-input v-model="costumer.Fields" type="textarea" readonly></b-input>
            </b-field>
            <b-field label="Стоимость" label-position="inside">
              <b-input v-model="costumer.Price"></b-input>
            </b-field>
            <b-field label="Период напоминания" label-position="inside">
              <b-select expanded v-model="costumer.Period">
                <option value="0">Отсутствует</option>
                <option value="1">Каждый день</option>
                <option value="2">Каждые 3 дня</option>
                <option value="3">Раз в неделю (понедельник)</option>
                <option value="4">Раз в 2 недели (понедельник)</option>
                <option value="5">Раз в месяц (первого числа)</option>
              </b-select>
            </b-field>
            <b-button
              type="is-info"
              @click="save(costumer.ID)"
              :disabled="loading || costumer.disabled"
            >Сохранить</b-button>
          </div>
        </div>
      </div>
    </div>
    <b-loading is-full-page :active="loading"></b-loading>
  </section>
</template>

<script>
import axios from 'axios'
import Navbar from '@/components/Navbar'
export default {
  name: 'costumers',
  components: { Navbar },
  data () {
    return {
      loading: false,
      costumers: []
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
        let response = await axios.get(this.$apiBase(this.username, this.password, 'costumers'))
        if (response.data.ok) {
          this.costumers = response.data.data
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
        for (let i = 0; i < this.costumers.length; i++) {
          if (this.costumers[i].ID === id) {
            response = await axios.post(
              this.$apiBase(this.username, this.password, 'costumers/modify?id=' + id),
              this.$qs.stringify(this.costumers[i])
            )
            break
          }
        }

        if (response.data.ok) {
          await this.load()
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
