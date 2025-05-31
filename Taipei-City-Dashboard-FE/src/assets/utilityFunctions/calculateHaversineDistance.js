/**
 * @typedef {Object} Coordinates
 * @property {number} latitude - The latitude of the coordinate.
 * @property {number} longitude - The longitude of the coordinate.
 */

/**
 * Calculates the Haversine distance between two points.
 *
 * @param {Coordinates} point1 - The first set of coordinates.
 * @param {Coordinates} point2 - The second set of coordinates.
 * @returns {number|null} The distance in kilometers between the two coordinates.
 */
export function calculateHaversineDistance(point1, point2) {
	try {
		if (
			isNaN(point1.latitude) ||
			isNaN(point1.longitude) ||
			isNaN(point2.latitude) ||
			isNaN(point2.longitude)
		) {
			throw new Error("All arguments must be numbers.");
		}

		const dx = point2.longitude - point1.longitude;
		const dy = point2.latitude - point1.latitude;
		
		return Math.sqrt(dx * dx + dy * dy) * 111;
	} catch {
		return null;
	}
}
