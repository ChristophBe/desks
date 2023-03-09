<template>


  <v-row>
    <v-col><h1>This Week</h1></v-col>


  </v-row>
  <v-row>
    <template v-for="(startOfDay, n) in days" :key="n">
      <v-col>
        <v-card :disabled="isDisabled(startOfDay)" :color="calculateType(startOfDay, n)"
                :variant="isDisabled(startOfDay)?'flat':'tonal'">
          <template v-slot:append>
            <v-btn icon="mdi-plus" elevation="0" @click="bookForDay(startOfDay)"></v-btn>
          </template>
          <template v-slot:title>
            {{ startOfDay.format("dddd") }}
          </template>
          <template v-slot:subtitle>
            {{ startOfDay.format("L") }}
          </template>


          <v-card-text>

            <v-chip>{{ n }} / 7 desks booked</v-chip>

          </v-card-text>
        </v-card>
      </v-col>
    </template>
  </v-row>

</template>
<script lang="ts">
import {defineComponent} from "vue";
import {mapActions, mapGetters, mapState} from "vuex";
import moment from "moment/moment";
import BookingUtils from "@/utils/booking-utils";
import Booking from "@/models/Booking";
import {Moment} from "moment";

interface DeskAvailabilityState {
  overlappingBookings: number | null
  bookingsOfColleagues: Booking[]
  days: Moment[]
}

export default defineComponent({
  name: 'desk-availability',
  components: {},
  data: (): DeskAvailabilityState => ({
    overlappingBookings: null,
    bookingsOfColleagues: [],
    days: []
  }),
  computed: {
    ...mapState("defaults", ["bookingDefaults", "bookingDefaultsLoading"]),
    ...mapGetters("bookings", ["getBookingsByRoomAndDay"]),
    ...mapState("bookings", ["bookings"]),
    ...mapState('user', ['user']),
  },
  mounted() {

    const startOfWeek = moment().startOf("week").add(1, "days")
    for (const i in [0, 0, 0, 0, 0]) {
      this.days.push(moment(startOfWeek).add((i), "days"))
    }
    this.fetchBookingDefaults()
    this.fetchBookings()
  },
  watch: {
    async bookingDefaults() {
      await this.fetchBookings()
      await this.loadAndCheck()
    },
    bookings() {
      this.loadAndCheck()
    }
  },
  methods: {

    ...mapActions("defaults", ["fetchBookingDefaults"]),
    ...mapActions("bookings", ["fetchBookingsByRoomAndDate"]),
    isDisabled: (startOfDay: Moment): boolean => startOfDay.isBefore(moment().startOf("day")),
    bookForDay: function (startOfDay: Moment) {
      const startDefault = moment(this.bookingDefaults.start);
      let start = startOfDay.startOf("day").add(startDefault.hour(), "hour").add(startDefault.minutes(), "minute")
      start = start.isBefore(moment.now()) ? BookingUtils.roundToNextFiveMinutes() : start

      const endDefault = moment(this.bookingDefaults.end);
      let end = moment().startOf("day").add(endDefault.hour(), "hour").add(endDefault.minutes(), "minute")
      end = end.isSameOrAfter(moment.now()) && end.isSameOrAfter(start) ? end : moment(start).add(1, "hour")
      const booking: Partial<Booking> = {
        room: this.bookingDefaults.room,
        start: start.toDate(),
        end: end.toDate()
      }
      this.$emit("book", booking)
    },

    async fetchBookings() {
      if (this.bookingDefaults && this.bookingDefaults.room) {
        await this.fetchBookingsByRoomAndDate({roomId: this.bookingDefaults.room.id, date: moment.now()})
      }
    },
    async loadAndCheck() {

      if (this.bookingDefaults?.room?.id) {

        this.checkOverlaps(this.getBookingsByRoomAndDay(this.bookingDefaults.room.id, moment.now()))
      }
    },

    checkOverlaps: function (bookings: Array<Booking>) {

      const start = moment().startOf("day")
      const end = moment().endOf("day")
      this.bookingsOfColleagues = bookings.filter(value => value.user.id !== this.user.id)
      const overlaps = BookingUtils.findOverlaps(bookings, start, end)
      this.overlappingBookings = new Set(overlaps.map(booking => booking.user.id)).size
    },

    calculateType(startOfDay: moment.Moment, number: number): string {
      if (this.isDisabled(startOfDay)) {
        return ""
      }
      if (number >= this.bookingDefaults.room.capacity) {
        return "red"
      }
      if (number > this.bookingDefaults.room.capacity * 0.8) {
        return "orange"
      }

      return "success"

    }
  }
});
</script>
