<template>

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
import {defineComponent, PropType} from "vue";
import {mapActions, mapGetters, mapState} from "vuex";
import moment from "moment/moment";
import BookingUtils from "@/utils/booking-utils";
import Booking from "@/models/Booking";
import {Moment, MomentInput} from "moment";
import AvailabilityCard from "@/components/booking-components/AvailabilityCard.vue"

export default defineComponent({
  name: 'weekOverviewWindow',
  props:{
    startOfWeek: {
      type: Object as PropType<MomentInput>,
      required: true
    },
  },
  components: {AvailabilityCard},

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
    },
  },
  methods: {

    ...mapActions("defaults", ["fetchBookingDefaults"]),
    ...mapActions("bookings", ["fetchBookingsByRoomAndDate"]),


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
      this.$emit("add-booking", booking)
    },

    async fetchBookings() {
      if (this.bookingDefaults && this.bookingDefaults.room) {
        await this.fetchBookingsByRoomAndDate({roomId: this.bookingDefaults.room.id, date: moment.now()})
      }
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
