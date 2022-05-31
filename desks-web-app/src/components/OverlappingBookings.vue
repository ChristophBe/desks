<template>
  <v-alert v-if="room"
           density="comfortable"
           :type="calculateType(overlappingBookings,room.capacity)"
           variant="contained-text"
           class="mb-2"
  >
    <div>{{ overlappingBookings }} / {{ room.capacity }} desks booked</div>
  </v-alert>
</template>


<script lang="ts">

import {defineComponent} from "vue";
import moment from "moment";
import {mapState} from "vuex";
import Booking from "@/models/Booking";

interface OverlappingBookingsData {
  overlappingBookings: number
  bookings: Array<Booking>
}

export default defineComponent({
  name: "OverlappingBookings",
  data(): OverlappingBookingsData {
    return {
      overlappingBookings: 0,
      bookings: [],
    }
  },
  // type inference enabled
  props: {
    room: {
      type: Object
    },
    date: String,
    start: String,
    end: String,
  },

  mounted() {
    this.loadAndCheck()
  },

  watch: {
    room() {
      this.loadAndCheck();
    },
    date() {
      this.loadAndCheck();
    },
    start() {
      this.checkOverlaps();
    },
    end() {
      this.checkOverlaps();
    }
  },

  computed: mapState('user', ['user']),

  methods: {

    formatTime(date: Date) {
      return moment(date).format("HH:mm")
    },

    async loadAndCheck() {
      console.log("load bookings for room and day")
      const res = await fetch(`/api/v1.0/rooms/${this.room?.id}/bookings?date=${moment(this.date).format('YYYY-MM-DD')}`)
      if (res.status >= 400) {
        return
      }
      this.bookings = await res.json()


      this.checkOverlaps()
    },

    checkOverlaps() {

      const start = moment(this.date).add(this.start)
      const end = moment(this.date).add(this.end)


      console.log(this.bookings)
      const overlaps = this.bookings.filter((booking) => {

        const bookingStart = moment(booking.start)
        const bookingEnd = moment(booking.end)
        return (start.isSameOrBefore(bookingStart) && end.isSameOrAfter(bookingEnd))
            || start.isBetween(bookingStart, bookingEnd)
            || end.isBetween(bookingStart, bookingEnd)
      })

      console.log(overlaps)
      this.overlappingBookings = overlaps.length
    },

    calculateType(numberOfOverlaps: number, capacity: number) {
      if (numberOfOverlaps >= capacity) {
        return "error"
      }
      if (numberOfOverlaps > capacity * 0.8) {
        return "warning"
      }

      return "success"

    }
  }

})
</script>

<style scoped>

</style>