<template>


  <v-row class="mt-4">
    <v-col>
      <h1 class="text-h4">{{ isPreviousWeekDisabled() ? "This Week" : "Week Overview" }}</h1>
      <h3 class="text-subtitle-1">{{$format.date(calculateNthDayOfWeek(1))}} - {{$format.date(calculateNthDayOfWeek(5))}}</h3>
    </v-col>

    <v-col align-self="center" class="d-flex justify-end">
      <v-btn icon variant="text" @click="()=>previousWeek()"
             :disabled="isPreviousWeekDisabled()">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-btn icon variant="text" @click="()=>nextWeek()">
        <v-icon>mdi-arrow-right</v-icon>
      </v-btn>
    </v-col>
  </v-row>
  <v-row>
    <template v-for="n in 5" :key="n">
      <v-col>
        <availability-card
            v-if="this.bookingDefaults"
            :startOfDay="calculateNthDayOfWeek(n)"
            :room="this.bookingDefaults.room"
            @add-booking="bookForDay(calculateNthDayOfWeek(n))"
        ></availability-card>
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
import AvailabilityCard from "@/components/booking-components/AvailabilityCard.vue"

interface DeskAvailabilityState {
  overlappingBookings: number | null
  bookingsOfColleagues: Booking[]
  startOfWeek: Moment
}

export default defineComponent({
  name: 'desk-availability',
  components: {AvailabilityCard},
  data: (): DeskAvailabilityState => ({
    overlappingBookings: null,
    bookingsOfColleagues: [],
    startOfWeek: moment().startOf("week")
  }),
  computed: {
    ...mapState("defaults", ["bookingDefaults", "bookingDefaultsLoading"]),
    ...mapGetters("bookings", ["getBookingsByRoomAndDay"]),
    ...mapState("bookings", ["bookings"]),
    ...mapState('user', ['user']),
  },
  mounted() {

    this.fetchBookingDefaults()
    this.fetchBookings()
  },
  watch: {
    async bookingDefaults() {
      await this.fetchBookings()
      await this.loadAndCheck()
    },
  },
  methods: {

    ...mapActions("defaults", ["fetchBookingDefaults"]),
    ...mapActions("bookings", ["fetchBookingsByRoomAndDate"]),

    nextWeek() {
      this.startOfWeek = moment(this.startOfWeek).add(1, "week")
    },
    previousWeek() {
      this.startOfWeek = moment(this.startOfWeek).subtract(1, "week")
    },
    isPreviousWeekDisabled(): boolean {
      return this.startOfWeek.isSameOrBefore(moment().startOf("week"))
    },
    isDisabled: (startOfDay: Moment): boolean => startOfDay.isBefore(moment().startOf("day")),

    calculateNthDayOfWeek(n: number): Moment {
      return moment(this.startOfWeek).add(n, "day")
    },

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
