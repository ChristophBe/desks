<template>
  <v-alert v-if="bookings.length <= 0">There are no other bookings for the office today.</v-alert>
  <v-list v-else>
    <v-list-item  class="px-0" two-line v-for="booking in bookings" :key="booking.id">
      <v-list-item-header>
        <v-list-item-title>{{ booking.user.givenName }} {{ booking.user.familyName }}</v-list-item-title>
        <v-list-item-subtitle>{{ formatTime(booking.start) }} - {{
            formatTime(booking.end)
          }}
        </v-list-item-subtitle>
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
}

export default defineComponent({
  name: "AlsoInTheRoom",
  data: (): alsoInTheRoomData => ({
    bookings: []
  }),
  // type inference enabled
  props: {
    roomId: {
      type: Number
    },
    date: String
  },

  async mounted() {

    const res = await fetch(`/api/v1.0/rooms/${this.roomId}/bookings?date=${moment(this.date).format('YYYY-MM-DD')}`)
    if (res.status >= 400) {
      return
    }
    const bookings = await res.json()
    this.bookings = bookings.filter((booking: Booking) => booking.user.id !== this.user.id)

  },

  computed: mapState('user', ['user']),

  methods: {

    formatTime(date: Date) {
      return moment(date).format("HH:mm")
    }
  }

})
</script>

<style scoped>

</style>