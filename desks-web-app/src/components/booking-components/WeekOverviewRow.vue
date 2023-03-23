<template>


  <v-row class="mt-4">
    <v-col>
      <h1 class="text-h4">Next Days Overview</h1>
      <h3 class="text-subtitle-1">{{$format.date(calculateNthNextDay(0))}} - {{$format.date(calculateLastWeekdayOfPeriod(calculateNthNextDay(6)))}} | {{bookingDefaults?.room?.name}}</h3>
    </v-col>

    <v-col align-self="center" class="d-flex justify-end flex-grow-0" :style="{width: 'fit-content'}">
      <v-btn icon variant="text"
             @click="()=>previousWeek()"
             :disabled="isPreviousWeekDisabled()">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-btn icon variant="text" @click="()=>nextWeek()">
        <v-icon>mdi-arrow-right</v-icon>
      </v-btn>
    </v-col>
  </v-row>
  <v-row>
    <v-col>
      <v-slide-group v-model="window" :style="{marginTop: '-1em'}">
        <v-slide-group-item v-for="(week,n) in weeks" :value="n" :key="n">
          <week-overview-window :start-of-week="week"  @add-booking="(booking) => $emit('book', booking)"></week-overview-window>
        </v-slide-group-item>
      </v-slide-group>
    </v-col>

  </v-row>

</template>
<script lang="ts">
import {defineComponent} from "vue";
import moment from "moment/moment";
import {Moment} from "moment";
import WeekOverviewWindow from "@/components/booking-components/WeekOverviewWindow.vue";
import {mapState} from "vuex";

interface DeskAvailabilityState {
  startOfWeek: Moment
  weeks: Moment[]
  window: number
}

export default defineComponent({
  name: 'desk-availability',
  components: {WeekOverviewWindow},
  data: (): DeskAvailabilityState => ({
    window: -1,
    weeks: [moment().startOf("day"), moment().isoWeekday(moment().startOf("day").isoWeekday() + 7).startOf("day")],
    startOfWeek: moment().startOf("day")
  }),
  computed:{
    ...mapState("defaults",["bookingDefaults"])
  },
  methods: {
    nextWeek() {
      const additionalWindow = (this.weeks[this.weeks.length-1]).isoWeekday((this.weeks[this.weeks.length-1]).isoWeekday() + 7).startOf("day")
      this.weeks = [...this.weeks,additionalWindow]
    },
    previousWeek() {
      this.startOfWeek = this.weeks[this.window]
    },
    isPreviousWeekDisabled(): boolean {
      return this.window <= 0
    },
    isDisabled: (startOfDay: Moment): boolean => startOfDay.isBefore(moment().startOf("day")),

    calculateNthNextDay(n: number): Moment {
      return moment(this.startOfWeek).add(n, "day")
    },
    calculateLastWeekdayOfPeriod(endOfPeriod: Moment): Moment {
      return endOfPeriod.isoWeekday() <= 5 ? endOfPeriod : endOfPeriod.isoWeekday(5);
    },

    getNext5Workingdays(start: Moment): Moment[] {
      const days = new Array<Moment>(5);
      days[0] = this.getNextWorkingDay(start);
      for(var i = 0; i < 5; i++) {
        days[i+1] = this.getNextWorkingDay(days[i]);
      }
      return days;
    },

    getNextWorkingDay(d: Moment) {
      const isoWeekday = d.isoWeekday();
      if(isoWeekday >= 5) {
        return d.isoWeekday(1);
      } else {
        return d.isoWeekday(isoWeekday + 1);
      }
    }
  }
});
</script>
