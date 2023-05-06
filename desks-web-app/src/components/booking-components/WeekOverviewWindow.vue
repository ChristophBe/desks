<template>

  <v-row>
    <template v-for="businessDay in getNext5BusinessDays()" :key="businessDay">
      <v-col class="d-flex flex-column justify-end" :class="{days: true, monday: isMonday(businessDay)}" v-if="this.bookingDefaults">
        <span id="dayNote" v-if="businessDay.startOf('day').isSame(today.startOf('day'), 'day')">Today</span>
        <span id="dayNote" v-else>Week {{businessDay.isoWeek()}}</span>
        <availability-card
            v-if="this.bookingDefaults"
            :startOfDay="businessDay"
            :room="this.bookingDefaults.room"
            @add-booking="bookForDay(businessDay)"
        ></availability-card>
      </v-col>
    </template>

  </v-row>

</template>
<style>

.days {
  padding-block: 1em;
}

.monday {
  border-left: 1px solid #585858;
}

.days:first-child {
  border-left: 0 solid #585858;
}


#dayNote {
  display: none;
  margin-bottom: 1em;
  color: #585858;
  line-height: 1;
  font-size: 1em;
}

.monday #dayNote, .days:first-child #dayNote {
  display: block;
}

</style>
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
    today: () => moment(),
  },
  async mounted() {
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
    ...mapActions("bookings", ["fetchBookingsByRoomAndTimespan"]),

    isDisabled: (startOfDay: Moment): boolean => startOfDay.isBefore(moment().startOf("day")),

    calculateNthNextDay(n: number): Moment {
      return moment(this.startOfWeek).add(n, "day")
    },

    isMonday(day: Moment) {
      return day.isoWeekday() === 1;
    },


    isWeekend(day: Moment) {
      return day.isoWeekday() > 5;
    },

    getNext5BusinessDays(): Moment[] {
      const days = new Array<Moment>(7);
      for(var i = 0; i < 7; i++) days[i] = this.calculateNthNextDay(i);
      return days.filter(d => !this.isWeekend(d));
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
        await this.fetchBookingsByRoomAndTimespan({roomId: this.bookingDefaults.room.id, from: this.startOfWeek, to: this.getNext5BusinessDays()[4]})
      }
    },
  }
});
</script>
