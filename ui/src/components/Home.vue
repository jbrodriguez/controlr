<template>
  <div>
    <section class="row">
      <div class="col-xs-12 center-xs">
        <h2>Scan this QR Code with the ControlR app to add your server</h2>
      </div>
    </section>

    <section class="row">
      <div v-if="hasBarcode" class="col-xs-12 center-xs">
        <img src="/qrcode.jpg" />
      </div>
      <div v-else />
    </section>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

import { GET_BARCODE } from '../constants'

export default Vue.extend({
  name: 'home',

  computed: {
    hasBarcode(): boolean {
      return this.$store.state.qrcode !== ''
    },
  },

  created() {
    this.$store.dispatch(GET_BARCODE)
  },

  //   methods: {
  // 	  this.image =
  //   },
})
</script>

<style scoped>
</style>
