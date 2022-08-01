<template>
  <v-progress-circular v-if="loading" indeterminate="" class="mx-4"></v-progress-circular>
  <v-alert v-else-if="bookings.length <= 0" class="mx-4">There are no other bookings for the same office.</v-alert>
  <v-list :density="density" v-else>
    <v-list-item two-line v-for="booking in bookings" :key="booking.id">
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
import {mapState} from "vuex";

interface alsoInTheRoomData {
  bookings: Array<Booking>
  loading: boolean
}

export default defineComponent({
  name: "AlsoInTheRoom",
  data: (): alsoInTheRoomData => ({
    bookings: [],
    loading: true
  }),
  // type inference enabled
  props: {
    roomId: {
      type: Number
    },
    density: String,
    date: String
  },

  async mounted() {
    const res = await fetch(`/api/v1.0/rooms/${this.roomId}/bookings?date=${moment(this.date).format('YYYY-MM-DD')}`)
    this.loading = false
    if (res.status >= 400) {
      return
    }
    const bookings = await res.json()
    this.bookings = bookings.filter((booking: Booking) => booking.user.id !== this.user.id)
  },
  computed: mapState('user', ['user']),
})
</script>

<style scoped>

</style>
