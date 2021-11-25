<template>
  <nav class="row pa1 middle-xs">
    <ul class="col-xs-12 col-sm-4">
      <li class="row start-xs">
        <div class="col-xs-12">
          <span>&nbsp;Copyright&nbsp;&copy;</span>
          <a class="dark-green" href="https://jbrio.net/posts/"
            >&nbsp;Juan B. Rodriguez</a
          >
        </div>
      </li>
    </ul>

    <ul class="col-xs-12 col-sm-3">
      <li class="middle-xs end-xs">
        <span>{{ version }}</span>
      </li>
    </ul>

    <ul class="col-xs-12 col-sm-5">
      <li class="row middle-xs end-xs">
        <div class="col-xs-12">
          <div class="flex end-xs">
            <a
              class="flex middle-xs ml2"
              title="Lime Technology"
              href="https://lime-technology.com"
              target="_blank"
              rel="noopener noreferrer"
            >
              <img
                src="../assets/unraid.svg"
                v-bind:style="{ height: iconSize }"
              />
            </a>

            <a
              class="flex middle-xs ml2"
              title="Juan B. Rodriguez"
              href="https://jbrio.net/posts/"
              target="_blank"
              rel="noopener noreferrer"
            >
              <img src="../assets/jbrio.png" />
            </a>
          </div>
        </div>
      </li>
    </ul>
  </nav>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'

export default Vue.extend({
  name: 'footer-bar',
  data() {
    return {
      iconSize: '45px',
    }
  },
  computed: mapState(['version']),
})
</script>

<style scoped>
</style>
