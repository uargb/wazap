<template>
  <section id="costumers">
    <navbar />
    <br/>
    <div class="container">
      <div class="columns">
        <div class="column is-12">
          <div class="box flex-row align-items-center">
            <p class="heading has-margin-right-15 has-margin-top-5">Управление</p>
            <b-button class="has-margin-right-10" type="is-info" icon-left="download" @click="download()">
              Экспорт клиентов в CSV
            </b-button>

            <b-button class="has-margin-right-10" type="is-info" icon-left="send" @click="mailing.active = true">
              Рассылка
            </b-button>
          </div>
        </div>
      </div>
      <div class="columns is-multiline">
        <div class="column is-4" v-for="costumer in costumers" :key="costumer.ID">
          <div class="box">
            <div class="buttons justify-space-between has-margin-bottom-20">
              <div></div>
              <b-tooltip label="Отключить клиента">
                <b-button type="is-danger"
                icon-left="delete"
                @click="remove(costumer.ID)"
                :loading="costumer.loading"
                :disabled="costumer.loading || costumer.disabled"></b-button>
              </b-tooltip>
            </div>
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
              <b-input v-model="costumer.Status"></b-input>
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
    <b-modal :active.sync="mailing.active" :width="640" scroll="keep" full-screen>
      <div class="modal-card justify-center has-padding-left-30 has-padding-right-30">
        <div>
          <p class="title">Рассылка карточки</p>
          <p class="subtitle">Позволяет отправить карточку сообщения всем определенному набору клиентов</p>
          <br>
          <b-field label="Статус получателей" label-position="inside">
            <b-select expanded v-model="mailing.status">
              <option value="all">Всем</option>
              <option v-for="status in statuses" :key="status" :value="status">{{status}}</option>
            </b-select>
          </b-field>
          <b-field label="ID карточки" label-position="inside">
            <b-input v-model="mailing.card"></b-input>
          </b-field>
          <br>
          <div class="buttons">
            <b-button type="button is-success" @click="send()">Отправить как можно скорее</b-button>
          </div>
        </div>
      </div>
    </b-modal>
    <b-loading is-full-page :active="loading"></b-loading>
  </section>
</template>

<script>
import axios from 'axios'
import Navbar from '@/components/Navbar'
var fileDownload = require('js-file-download')
export default {
  name: 'costumers',
  components: { Navbar },
  data () {
    return {
      loading: false,
      costumers: [],
      statuses: [],
      mailing: {
        active: false,
        status: '',
        cardId: 0
      }
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
        let response = await axios.get(this.$apiBase(this.username, this.password, 'costumers'))
        if (response.data.ok) {
          this.costumers = response.data.data
          this.costumers.forEach(costumer => {
            if (!this.statuses.includes(costumer.Status)) {
              this.statuses.push(costumer.Status)
            }
          })
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
    },
    async remove (id) {
      this.loading = true
      try {
        let response = await axios.post(this.$apiBase(this.username, this.password, 'costumers/remove?id=' + id))
        if (response.data.ok) {
          for (let i = 0; i < this.costumers.length; i++) {
            if (this.costumers[i].ID === id) {
              this.costumers[i].disabled = true
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
    async download () {
      this.loading = true
      try {
        let response = await axios.get(this.$apiBase(this.username, this.password, 'costumers/export'))
        fileDownload(response.data, 'costumers.csv')
      } catch (error) {
        this.$error(this, error.message)
      } finally {
        this.loading = false
      }
    },
    async send () {
      this.loading = true
      try {
        let response = await axios.post(
          this.$apiBase(this.username, this.password, 'costumers/send'),
          this.$qs.stringify(this.mailing)
        )
        this.mailing.active = false
      } catch (error) {
        this.$error(this, error.message)
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
