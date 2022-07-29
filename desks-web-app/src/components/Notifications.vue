<template>
  <v-snackbar
      v-model="snackbar"
      :color="notificationType === 'error' ? 'error': 'primary'"
      :timeout="timeout"
  >
    {{ text }}

    <template v-slot:actions>
      <v-btn
          variant="text"
          @click="snackbar = false"
      >
        Close
      </v-btn>
    </template>
  </v-snackbar>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {mapState} from "vuex";

export default defineComponent({
  name: "app-notifications",
  data: () => ({
    snackbar: false,
    timeout: 2000
  }),
  computed:{
    ... mapState("notification", ["text","notificationType", "time"])
  },
  watch:{
    time (newText, oldText) {

      console.log(newText,oldText)
      this.snackbar = false;
      this.snackbar = newText !== ""
    }
  }
});

</script>

<style scoped>

</style>
