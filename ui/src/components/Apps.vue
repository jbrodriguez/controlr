<template>
  <div>
    <section class="row middle-xs mb2">
      <div class="col-xs-12 start-xs middle-xs ttu mb2">
        <span class="f3">Apps</span>
      </div>
    </section>
    <section class="row middle-xs mb2">
      <div class="col-xs-12 start-xs middle-xs mb2">
        <div>
          <strong class="mr2">READ:</strong>
          <span>User can see the docker/vm, but actions are disabled</span>
          <br />
          <strong class="mr2">WRITE:</strong>
          <span>User can edit/remove the docker/vm</span>
          <br />
          <strong class="mr2">EXEC:</strong>
          <span>User can start/stop the docker/vm</span>
          <p class="mt2">
            Choosing any of the permissions will make the app visible in the
            ControlR app
          </p>
        </div>
      </div>
    </section>
    <section class="row middle-xs mb2">
      <div class="col-xs-12 start-xs middle-xs">
        <div v-if="!hasApps">No dockers/vms are available</div>
        <table v-else>
          <thead>
            <tr>
              <th class="tl w500">APP</th>
              <th class="tc w150">READ</th>
              <th class="tc w150">WRITE</th>
              <th class="tc w150">EXEC</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="app in apps" :key="app.name">
              <td class="flex flex-row items-center">
                <div>
                  <img :src="app.logo" alt="icon" class="h2" />
                </div>
                <div class="ml2" />
                <div>{{ app.name }}</div>
              </td>
              <td class="tc">
                <div class="flex' items-center justify-center">
                  <input
                    type="checkbox"
                    :checked="app.read"
                    @click="setPerm({ user, name: app.name, perm: 'read' })"
                  />
                </div>
              </td>
              <td class="tc">
                <div class="flex items-center justify-center">
                  <input
                    type="checkbox"
                    :checked="app.write"
                    @click="setPerm({ user, name: app.name, perm: 'write' })"
                  />
                </div>
              </td>
              <td class="tc">
                <div class="flex items-center justify-center">
                  <input
                    type="checkbox"
                    :checked="app.exec"
                    @click="setPerm({ user, name: app.name, perm: 'exec' })"
                  />
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import Button from './Button.vue'

import { SET_PERM } from '../constants'
import { IUser, IState, IApp, ISetPermArgs } from '../types'

export default Vue.extend({
  name: 'apps',

  components: { Button },

  props: ['user'],

  computed: {
    apps(): IApp[] {
      if (this.user === '') {
        return []
      }

      const state: IState = this.$store.state
      return state.appOrder.map(name => state.apps[this.user][name])
    },

    hasApps(): boolean {
      return this.$store.state.appOrder.length > 0
    },
  },

  methods: {
    setPerm(args: ISetPermArgs) {
      this.$store.commit(SET_PERM, args)
    },
  },
})
</script>

<style scoped>
.w500 {
  width: 500px;
}

.w150 {
  width: 150px;
}
</style>
