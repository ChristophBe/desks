<template>
  <v-list>
    <v-list-item two-line v-for="booking in bookingsSorted" :key="booking.id">
      <v-list-item-avatar start>
        <v-avatar color="primary">
          {{ booking.user.givenName.charAt(0) }}{{ booking.user.familyName.charAt(0) }}
        </v-avatar>


      </v-list-item-avatar>
      <v-list-item-header>
        <v-list-item-title>{{ booking.user.givenName }} {{ booking.user.familyName }}</v-list-item-title>
        <v-list-item-subtitle>{{ $format.timeRange(booking.start, booking.end) }}</v-list-item-subtitle>
      </v-list-item-header>

    </v-list-item>
  </v-list>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Booking from "@/models/Booking";
import moment from "moment";
import {th} from "@vuetify/nightly/locale";

interface BookingListState {
  bookingsSorted: Booking[]
}

export default defineComponent({
  name: "BookingList",
  data: (): BookingListState => ({
    bookingsSorted: []
  }),
  props: {
    bookings: {
      type: Array
    }
  },

  watch: {
    bookings() {
      this.bookingsSorted = this.bookings ? [...this.bookings].sort(this.compareByName) : [];

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