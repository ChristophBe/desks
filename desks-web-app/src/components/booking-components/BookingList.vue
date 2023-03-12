<template>
  <v-list>
    <v-list-item two-line v-for="booking in [...bookings].sort(compareByName)" :key="booking.id">
      <template v-slot:prepend>
        <v-avatar color="primary">
          {{ booking.user.givenName.charAt(0) }}{{ booking.user.familyName.charAt(0) }}
        </v-avatar>
      </template>
      <template v-slot:title>{{ booking.user.givenName }} {{ booking.user.familyName }}</template>
      <template v-slot:subtitle>{{ $format.timeRange(booking.start, booking.end) }}</template>
    </v-list-item>
  </v-list>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Booking from "@/models/Booking";
import moment from "moment";


export default defineComponent({
  name: "BookingList",

  props: {
    bookings: {
      type: Array
    }
  },

  methods: {
    compareByName(a: Booking, b: Booking): number {
      let comparison = a.user.familyName.localeCompare(b.user.familyName)
      if (comparison === 0) {
        comparison = a.user.givenName.localeCompare(b.user.givenName)
      }
      if (comparison === 0) {
        comparison = moment(a.start).diff(moment(b.start))
      }
      return comparison
    },
  }

})
</script>

<style scoped>

</style>
