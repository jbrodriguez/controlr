<template>
  <div class="container flex flex-column body">
    <header class="mv3">
      <header-bar></header-bar>
    </header>
    <main>
      <div class="flex flex-column">
        <feedback v-if="hasError" />
        <router-view></router-view>
      </div>
    </main>
    <footer class="mv3 bg-light-gray">
      <footer-bar></footer-bar>
    </footer>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'

import HeaderBar from './components/HeaderBar.vue'
import FooterBar from './components/FooterBar.vue'
import Feedback from './components/Feedback.vue'

export default Vue.extend({
  name: 'app',

  components: { HeaderBar, FooterBar, Feedback },

  computed: mapState(['hasError']),
})
</script>

<style>
* {
  box-sizing: border-box;
}

html {
  -ms-text-size-adjust: 100%;
  -webkit-text-size-adjust: 100%;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

body {
  margin: 0;
  font-size: 18px;
}

h1,
h2,
h3,
h4,
p,
blockquote,
figure,
ol,
ul,
span,
pre {
  margin: 0;
  padding: 0;
}

h1,
h2,
h3,
h4 {
  font-size: inherit;
}

a {
  text-decoration: none;
  color: inherit;
  -webkit-transition: 0.3s;
  -moz-transition: 0.3s;
  transition: 0.3s;
}

img {
  max-width: 100%;
  height: auto;
  border: 0;
}

input {
  border: none;
}

body {
  font: 100%/1.625em -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto',
    'Lucida Sans Unicode', 'Lucida Grande', sans-serif;
  color: #001b44;
  background: #f4f4f4;
}

h1,
h2,
h3 {
  font-weight: bold;
  line-height: 1.2;
}

h1 {
  font-size: 2em;
  margin-bottom: 0.17em;
}

h2 {
  font-size: 1.5em;
  margin-bottom: 0.33em;
}

h3 {
  font-size: 1.17em;
  margin-bottom: 0.5em;
}

ul {
  list-style-position: inside;
}

ul li {
  list-style-type: none;
}

ol {
  list-style-position: inside;
}

ol li {
  list-style-type: none;
}

table {
  border-collapse: collapse;
  border-spacing: 0;
}

/* global */
.body {
  min-height: 100vh;
  margin-right: auto;
  margin-left: auto;
}

main {
  flex: 1;
}
</style>