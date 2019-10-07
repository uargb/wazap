<template>
  <section id="navbar">
    <b-navbar type="is-primary">
      <template slot="brand">
        <b-navbar-item tag="p">
          <img src="@/assets/logo.png"/>
        </b-navbar-item>
        <b-navbar-item tag="span">Панель менеджера</b-navbar-item>
      </template>
      <template slot="start">
        <b-navbar-item :active="route === 'general'"
          tag="router-link"
          to="general">Основная информация</b-navbar-item>
        <b-navbar-item :active="route === 'qa'"
          tag="router-link"
          to="qa">Меню бота</b-navbar-item>
        <b-navbar-item :active="route === 'costumers'"
          tag="router-link"
          to="costumers">Клиенты</b-navbar-item>
      </template>

      <template slot="end">
        <b-navbar-item tag="div" class="has-margin-right-25">
          <b-button type="button is-primary" @click="isFileModalActive = true">Загрузить файл(ы)</b-button>
        </b-navbar-item>
        <b-navbar-item tag="p" class="is-text-light">
          <b-icon icon="account"></b-icon>
          <span>{{this.username}}</span>
        </b-navbar-item>
        <b-navbar-item tag="div">
          <router-link class="button is-primary"
            to="/"
            >Выход</router-link>
        </b-navbar-item>
      </template>
    </b-navbar>

    <b-modal :active.sync="isFileModalActive" :width="640" scroll="keep" full-screen>
      <div class="modal-card justify-center has-padding-left-30 has-padding-right-30">
        <div>
          <p class="title">Загрузка файлов</p>
          <p class="subtitle">Позволяет загрузить один или несколько файлов на сервер для использования ботом</p>
          <p>
            Доступные форматы: картинки - jpg, png, gif | видео - mp4 | вложения - pdf | импорт клиентов - csv
          </p>
          <br>
          <b-field>
            <b-upload v-model="dropFiles" multiple drag-drop>
              <section class="section">
                <div class="content has-text-centered">
                  <p><b-icon icon="upload" size="is-large"></b-icon></p>
                  <p>Перетащите сюда файлы или нажмите для выбора</p>
                </div>
              </section>
            </b-upload>
          </b-field>
          <div class="tags">
            <span v-for="(file, index) in dropFiles" :key="index" class="tag is-medium is-primary">
              {{file.name}}
              <button class="delete is-small" type="button" @click="deleteDropFile(index)" />
            </span>
          </div>
          <br>
          <div class="buttons">
            <b-button type="button is-success" @click="uploadFiles()"
            :loadging="fileUploading" :disabled="fileUploading">Выгрузить</b-button>
            <b-button type="button is-light" @click="cancelUploading()"
            :loadging="fileUploading" :disabled="fileUploading">Отмена</b-button>
          </div>
        </div>
      </div>
    </b-modal>
  </section>
</template>

<script>
import axios from 'axios'
export default {
  name: 'navbar',
  data () {
    return {
      fileUploading: false,
      isFileModalActive: false,
      dropFiles: []
    }
  },
  computed: {
    route () { return this.$route.name },
    username () { return this.$store.state.username },
    password () { return this.$store.state.password }
  },
  methods: {
    deleteDropFile (index) {
      this.dropFiles.splice(index, 1)
    },
    async uploadFiles () {
      this.fileUploading = true
      try {
        this.dropFiles.forEach(async file => {
          let formData = new FormData()
          formData.append('file', file)
          let response = await axios.post(
            this.$apiBase(this.username, this.password, 'upload'),
            formData,
            {
              headers: { 'Content-Type': 'multipart/form-data' }
            }
          )
          if (!response.data.ok) {
            this.$error(this, response.data.message)
          }
        })
      } catch (error) {
        this.$error(this, error.message)
      } finally {
        this.fileUploading = false
        this.cancelUploading()
      }
    },
    cancelUploading () {
      this.dropFiles = []
      this.isFileModalActive = false
    }
  }
}
</script>
