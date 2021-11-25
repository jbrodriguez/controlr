<template>
  <nav class="row">
    <ul class="col-xs-12 col-sm-2">
      <li class="flex bg-light-gray center-xs pa1">
        <img src="../assets/controlr.png" />
      </li>
    </ul>

    <ul class="col-xs-12 col-sm-10">
      <div class="bg-light-gray pa1">
        <li class="row middle-xs">
          <div class="col-xs-12 col-sm-10">
            <a
              title="Download on the App Store"
              href="https://itunes.apple.com/us/app/controlr/id1152586217?ls=1&mt=8"
              target="_blank"
              rel="noopener noreferrer"
              class="mr3"
            >
              <img src="../assets/appstore.png" />
            </a>

            <a
              title="Get it on Google Play"
              href="https://play.google.com/store/apps/details?id=io.jbrodriguez.controlr&utm_source=global_co&utm_medium=prtnr&utm_content=Mar2515&utm_campaign=PartBadge&pcampaignid=MKT-Other-global-all-co-prtnr-py-PartBadge-Mar2515-1"
              target="_blank"
              rel="noopener noreferrer"
            >
              <img src="../assets/gbadge.png" />
            </a>
          </div>

          <div class="col-xs-12 col-sm-1">
            <div v-if="isBusy" class="loader"></div>
          </div>

          <div class="col-xs-12 col-sm-1">
            <div class="row middle-xs end-xs">
              <div class="col-xs-12">
                <div class="flex middle-xs end-xs">
                  <a
                    class="flex"
                    href="https://twitter.com/jbrodriguezio"
                    title="@jbrodriguezio"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <img class="c-icon" src="../assets/twitter.svg" />
                  </a>

                  <a
                    class="flex ml2"
                    href="https://github.com/jbrodriguez"
                    title="github.com/jbrodriguez"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <img class="c-icon" src="../assets/github.svg" />
                  </a>
                  <img
                    src="../assets/v.png"
                    v-bind:style="{ height: iconSize }"
                  />
                </div>
              </div>
            </div>
          </div>
        </li>
      </div>
    </ul>
  </nav>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'

import Welcome from './Welcome.vue'
import Login from './Login.vue'

export default Vue.extend({
  name: 'header-bar',

  data() {
    return {
      iconSize: '45px',
    }
  },

  components: { Welcome, Login },

  computed: mapState(['isAuthorized', 'isBusy']),
})
</script>

<style scoped>
.c-icon {
  height: 25px;
  width: 25px;
  fill: red;
}

.loader {
  animation: spin 1s linear infinite;
  border: 3px solid #ddd;
  border-top: 3px solid red;
  border-radius: 50%;
  height: 25px;
  width: 25px;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
