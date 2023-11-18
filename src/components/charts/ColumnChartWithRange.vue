<!-- Developed by  -->

<script setup>
import { onMounted, ref } from "vue";
import { useMapStore } from "../../store/mapStore";

const props = defineProps([
	"chart_config",
	"activeChart",
	"series",
	"map_config",
]);
const newSeries = ref([]);
const mapStore = useMapStore();

const chartOptions = ref({
	chart: {
		stacked: true,
		toolbar: {
			show: false,
		},
	},
	colors: props.chart_config.color,
	dataLabels: {
		enabled: props.chart_config.categories ? false : true,
		offsetY: 20,
	},
	grid: {
		show: false,
	},
	legend: {
		show: props.chart_config.categories ? true : false,
	},
	plotOptions: {
		bar: {
			borderRadius: 5,
		},
	},
	stroke: {
		colors: ["#282a2c"],
		show: true,
		width: 2,
	},
	tooltip: {
		// The class "chart-tooltip" could be edited in /assets/styles/chartStyles.css
		custom: function ({ series, seriesIndex, dataPointIndex, w }) {
			return (
				'<div class="chart-tooltip">' +
				"<h6>" +
				w.globals.labels[dataPointIndex] +
				`${
					props.chart_config.categories
						? "-" + w.globals.seriesNames[seriesIndex]
						: ""
				}` +
				"</h6>" +
				"<span>" +
				series[seriesIndex][dataPointIndex] +
				` ${props.chart_config.unit}` +
				"</span>" +
				"</div>"
			);
		},
	},
	xaxis: {
		axisBorder: {
			show: false,
		},
		axisTicks: {
			show: false,
		},
		categories: props.chart_config.categories
			? props.chart_config.categories
			: [],
		labels: {
			offsetY: 5,
		},
		type: "category",
	},
});
function handleOnChange(series, input) {
	newSeries.value = [];
	let newData = [];
	for (let i = 0; i < series.length; i++) {
		newData.push(series[i].data[input]);
	}
	newSeries.value.push({ name: "某年", data: newData });
}

function setDefault(series) {
	let defaultData = [];
	let defaultS = [];

	for (let i = 0; i < series.length; i++) {
		defaultData.push(series[i].data[0]);
	}
	defaultS.push({ name: "某年", data: defaultData });
	return defaultS;
}

const selectedIndex = ref(null);

function handleDataSelection(e, chartContext, config) {
	if (!props.chart_config.map_filter) {
		return;
	}
	const toFilter = config.dataPointIndex;
	if (toFilter !== selectedIndex.value) {
		mapStore.addLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`,
			props.chart_config.map_filter[0],
			props.chart_config.map_filter[1][toFilter]
		);
		selectedIndex.value = toFilter;
	} else {
		mapStore.clearLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`
		);
		selectedIndex.value = null;
	}
}
</script>

<template>
	<div v-if="activeChart === 'ColumnChartWithRange'">
		<apexchart
			width="100%"
			height="200px"
			type="bar"
			:options="chartOptions"
			:series="newSeries.length ? newSeries : setDefault(series)"
			@dataPointSelection="handleDataSelection"
		></apexchart>
		<div>
			<input
				type="range"
				v-model="sliderValue"
				:min="0"
				:max="series[0].data.length - 1"
				@change="handleOnChange(series, sliderValue)"
			/>
			<p>
				{{
					isNaN(sliderValue)
						? "103年"
						: (103 + parseInt(sliderValue)).toString() + "年"
				}}
			</p>
		</div>
	</div>
</template>
