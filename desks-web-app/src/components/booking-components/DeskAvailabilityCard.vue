<template>
  <v-card v-if="!bookingDefaultsLoading && bookingDefaults != null">
    <v-card-item>
      <div class="d-flex justify-space-between align-center">
        <div>
          <v-card-title>Desk Availability</v-card-title>
          <v-card-subtitle>{{ bookingDefaults.room.name }}</v-card-subtitle>
        </div>
        <div v-if="overlappingBookings !== null">

            <v-menu open-on-hover :disabled="bookingsOfColleagues.length<=0">
              <template v-slot:activator="{ props }">
                <v-chip v-bind="props" :color="calculateType()">{{
                    overlappingBookings
                  }} /
                  {{ bookingDefaults.room.capacity }} desks booked today
                </v-chip>
              </template>
                <booking-list v-if="bookingsOfColleagues.length>0" :bookings="bookingsOfColleagues"></booking-list>

            </v-menu>
        </div>
      </div>
    </v-card-item>
    <v-card-actions>
      <v-btn variant="text" @click="booksForToday">
        book now
      </v-btn>
    </v-card-actions>
  </v-card>
</template>
<script lang="ts">
import {defineComponent} from "vue";
import {mapActions, mapGetters, mapState} from "vuex";
import moment from "moment/moment";
import BookingUtils from "@/utils/booking-utils";
import Booking from "@/models/Booking";
import BookingList from "@/components/booking-components/BookingList.vue";

interface DeskAvailabilityState {
  overlappingBookings: number | null
  bookingsOfColleagues: Booking[]
}

export default defineComponent({
  name: 'desk-availability',
  components: {BookingList},
  data: (): DeskAvailabilityState => ({
    overlappingBookings: null,
    bookingsOfColleagues: []
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
    bookings() {
      this.loadAndCheck()
    }
  },
  methods: {

    ...mapActions("defaults", ["fetchBookingDefaults"]),
    ...mapActions("bookings", ["fetchBookingsByRoomAndDate"]),

    booksForToday: function () {
      const startDefault = moment(this.bookingDefaults.start);
      let start = moment().startOf("day").add(startDefault.hour(), "hour").add(startDefault.minutes(), "minute")
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
