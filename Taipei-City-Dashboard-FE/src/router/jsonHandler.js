import axios from "axios";

export default {
	data() {
		return {
			jsonData: null,
			modifiedData: null,
		};
	},
	methods: {
		async uploadData(updateData) {
			try {
				const response = await axios.get("/mapData/incident.geojson");
				this.jsonData = response.data;
				this.modifyData(updateData);
				this.exportData();
			} catch (error) {
				console.error("error fetching data...");
			}
		},
		modifyData(updateData) {
			// Deep copy of jsonData
			this.modifiedData = JSON.parse(JSON.stringify(this.jsonData));
			// Modify the deep copied data
			this.modifiedData.features.push(updateData);
		},
		async exportData() {
			try {
				console.log(this.modifiedData);
				const res = await axios.put(
					"/mapData/incident.geojson",
					JSON.stringify(this.modifiedData)
				);
				console.log(res);
				console.log("data export successfully");
			} catch (error) {
				console.log("error exporting data...");
			}
		},
	},
};
