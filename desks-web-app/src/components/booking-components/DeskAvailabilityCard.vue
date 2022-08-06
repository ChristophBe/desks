<template>
  <v-card v-if="!bookingDefaultsLoading && bookingDefaults != null">
    <v-card-item>
      <div class="d-flex justify-space-between align-center">
        <div>
          <v-card-title>Desk Availability</v-card-title>
          <v-card-subtitle>{{ bookingDefaults.room.name }}</v-card-subtitle>
        </div>
        <v-fade-transition>
          <v-chip v-if="overlappingBookings !== null" :color="calculateType()">{{ overlappingBookings }} /
            {{ bookingDefaults.room.capacity }} desks booked today
          </v-chip>
        </v-fade-transition>

      </div>
    </v-card-item>
    <v-card-actions>
      <v-btn variant="text" @click="booksForToday">
        Book desk
      </v-btn>
    </v-card-actions>
  </v-card>
</template>
<script lang="ts">
import {defineComponent} from "vue";
import {mapActions, mapState} from "vuex";
import moment from "moment/moment";
import BookingUtils from "@/utils/booking-utils";
import Booking from "@/models/Booking";

interface DeskAvailabilityState {
  overlappingBookings: number | null
}

export default defineComponent({
  name: 'desk-availability',
  data: (): DeskAvailabilityState => ({
    overlappingBookings: null
  }),
  computed: {
    ...mapState("defaults", ["bookingDefaults", "bookingDefaultsLoading"])
  },
  mounted() {
    this.loadAndCheck()
    this.fetchBookingDefaults()
  },
  watch: {
    bookingDefaults() {
      this.loadAndCheck()
    }
  },
  methods: {
    ...mapActions("defaults", ["fetchBookingDefaults"]),

    booksForToday: function () {
      const startDefault = moment(this.bookingDefaults.start);
      let start = moment().startOf("day").add(startDefault.hour(), "hour").add(startDefault.minutes(), "minute")
      start = start.isBefore(moment.now()) ? this.roundToNextFiveMinutes() : start

      const endDefault = moment(this.bookingDefaults.start);
      let end = moment().startOf("day").add(endDefault.hour(), "hour").add(endDefault.minutes(), "minute")
      end = end.isSameOrAfter(moment.now()) && end.isSameOrAfter(start) ? end : moment(start).add(1, "hour")
      const booking: Partial<Booking> = {
        room: this.bookingDefaults.room,
        start: start.toDate(),
        end: end.toDate()
      }
      this.$emit("book", booking)
    },
    roundToNextFiveMinutes() {
      let now = moment().startOf("hour")
      let min = (5 * Math.floor(moment().diff(now, "minute") / 5)) + 5
      return now.add(min, "minutes")
    },

    async loadAndCheck() {

      if (!this.bookingDefaults) {
        return
      }
      console.log("load bookings for room and day")
      const res = await fetch(`/api/v1.0/rooms/${this.bookingDefaults.room?.id}/bookings?date=${moment().format('YYYY-MM-DD')}`)
      if (res.status >= 400) {
        return
      }
      const bookings = await res.json()


      this.checkOverlaps(bookings)
    },

    checkOverlaps: function (bookings: Array<Booking>) {

      const start = moment().startOf("day")
      const end = moment().endOf("day")

      console.log("test", start, end)
      const overlaps = BookingUtils.findOverlaps(bookings, start, end)
      this.overlappingBookings = new Set(overlaps.map(booking => booking.user.id)).size
    },

    calculateType() {

      if (this.overlappingBookings === null) {
        return "error"
      }
      if (this.overlappingBookings >= this.bookingDefaults.room.capacity) {
        return "error"
      }
      if (this.overlappingBookings > this.bookingDefaults.room.capacity * 0.8) {
        return "warning"
      }

      return "success"

    }
  }
});
</script>
