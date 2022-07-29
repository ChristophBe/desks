<template>
  <v-btn
      color="orange-darken-2"
      @click="snackbar = true"
  >
    Open Snackbar
  </v-btn>

  <v-snackbar
      v-model="snackbar"
      :timeout="timeout"
  >
    {{ text }}

    <template v-slot:actions>
      <v-btn
          color="blue"
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
import OverlappingBookings from "@/components/booking-components/OverlappingBookings.vue";
import moment from "moment/moment";

export default defineComponent({
      name: "Info",
      data: () => ({
        valid: false,
        date: moment().add(1, 'days').format("YYYY-MM-DD"),
        room: null,
        start: '09:00',
        end: '17:00',
        hasOwnOverlaps: false
      })
    );

</script>

<style scoped>

</style>
