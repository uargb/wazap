<template>
  <b-navbar type="is-primary">
    <template slot="brand">
      <b-navbar-item tag="p">
        <img src="@/assets/logo.png"/>
      </b-navbar-item>
      <b-navbar-item tag="span">Административный кабинет менеджера</b-navbar-item>
    </template>
    <template slot="start">
      <b-navbar-item :active="route === 'general' ? true : false"
        tag="router-link"
        to="general">Основная информация</b-navbar-item>
      <b-navbar-item :active="route === 'qa' ? true : false"
        tag="router-link"
        to="qa">Меню бота</b-navbar-item>
      <!-- <b-navbar-item :active="route === 'plans' ? true : false"
        tag="router-link"
        to="plans">Тарифы</b-navbar-item> -->
      <b-navbar-item :active="route === 'costumers' ? true : false"
        tag="router-link"
        to="costumers">Клиенты</b-navbar-item>
    </template>

    <template slot="end">
      <b-navbar-item tag="p" class="is-text-light">
        <b-icon icon="account"></b-icon>
        <span>{{this.username}}</span>
      </b-navbar-item>
      <b-navbar-item tag="div">
        <router-link class="button is-link"
          to="/"
          >Выход</router-link>
      </b-navbar-item>
    </template>
  </b-navbar>
</template>

<script>
export default {
  name: 'navbar',
  computed: {
    route () {
      return this.$route.name
    },
    username () {
      return this.$store.state.username
    }
  }
}
</script>
