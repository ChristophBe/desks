<template>


  <v-container class="container align-center justify-center d-flex">

    <v-card class="d-flex align-center flex-column" variant="plain">
      <v-card-text>
        <v-progress-circular indeterminate="" size="large"></v-progress-circular>

      </v-card-text>
      <v-card-item>
        <v-card-subtitle>
          {{ text ?? "Loading..." }}
        </v-card-subtitle>
      </v-card-item>

    </v-card>


  </v-container>
</template>

<script lang="ts">
import {defineComponent} from "vue";

export default defineComponent({
  name: "app-loading",
  props: {
    text: String
  }
})
</script>

<style scoped>
.container {
  height: 100%;
}
</style>