<template>
  <v-expand-transition>
    <v-alert v-if="room && date"
             density="comfortable"
             :type="calculateType(overlappingBookings,room.capacity)"
             variant="text"
             class="mb-2"
    >
      <div>{{ overlappingBookings }} / {{ room.capacity }} desks booked</div>
    </v-alert>
  </v-expand-transition>

</template>


<script lang="ts">

import {defineComponent} from "vue";
import moment from "moment";
import {mapActions, mapGetters, mapState} from "vuex";
import Booking from "@/models/Booking";
import BookingUtils from "@/utils/booking-utils";

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

  computed: {
    ...mapState('user', ['user']),
    ...mapGetters('bookings', ['getBookingsByRoomAndDay'])
  },

  methods: {
    ...mapActions("bookings", ["fetchBookingsByRoomAndDate"]),
    async loadAndCheck() {

      if (!this.room || !this.date) {
        return
      }
      await this.fetchBookingsByRoomAndDate({roomId: this.room.id, date: this.date})
      this.bookings = this.getBookingsByRoomAndDay(this.room.id, this.date)

      this.checkOverlaps()
    },

    checkOverlaps: function () {

      const start = moment(this.date).add(this.start)
      const end = moment(this.date).add(this.end)
      const overlaps = BookingUtils.findOverlaps(this.bookings, start, end)
      this.overlappingBookings = overlaps.length
    },

    calculateColor(numberOfOverlaps: number, capacity: number) {
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
