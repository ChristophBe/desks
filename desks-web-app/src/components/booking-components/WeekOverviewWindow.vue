<template>

  <v-row>
    <template v-for="n in [0,1,2,3,4,5,6]" :key="n">
      <v-col v-if="this.bookingDefaults && calculateNthNextDay(n).isoWeekday() <= 5">
        <availability-card
            v-if="this.bookingDefaults"
            :startOfDay="calculateNthNextDay(n)"
            :room="this.bookingDefaults.room"
            @add-booking="bookForDay(calculateNthNextDay(n))"
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
  async mounted() {

    console.log(this.startOfWeek);

    await this.fetchBookingDefaults()
    await this.fetchBookings()
  },
  watch: {
    async bookingDefaults() {
      await this.fetchBookings()
    },
  },
  methods: {

    ...mapActions("defaults", ["fetchBookingDefaults"]),
    ...mapActions("bookings", ["fetchBookingsByRoomAndWeek"]),

    isDisabled: (startOfDay: Moment): boolean => startOfDay.isBefore(moment().startOf("day")),

    calculateNthNextDay(n: number): Moment {
      console.log(moment(this.startOfWeek).add(n, "day"))
      return moment(this.startOfWeek).add(n, "day")
    },

    setDay(base:MomentInput, day: Moment): Moment{
      return moment(base).set({'year': day.year(), 'month': day.month(),'day': day.day()});
    },

    bookForDay(startOfDay: Moment) {

      let start = this.setDay(this.bookingDefaults.start, startOfDay);
      start = start.isBefore(moment.now()) ? BookingUtils.roundToNextFiveMinutes() : start

      let end = this.setDay(this.bookingDefaults.end, startOfDay)
      end = end.isSameOrAfter(start) ? end : moment(start).add(1, "hour")
      const booking: Partial<Booking> = {
        room: this.bookingDefaults.room,
        start: start.toDate(),
        end: end.toDate()
      }
      this.$emit("add-booking", booking)
    },

    async mounted(){
      await this.fetchBookings()
    },

    async fetchBookings() {
      if (this.bookingDefaults && this.bookingDefaults.room) {
        await this.fetchBookingsByRoomAndWeek({roomId: this.bookingDefaults.room.id, date: this.startOfWeek})
      }
    },
  }
});
</script>
