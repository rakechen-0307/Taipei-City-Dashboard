import pandas as pd
import json
import sys
import numpy as np
from typing import Dict, List, Any

# Custom JSON encoder to handle NumPy types
class NumpyEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, np.integer):
            return int(obj)
        elif isinstance(obj, np.floating):
            return float(obj)
        elif isinstance(obj, np.ndarray):
            return obj.tolist()
        return super(NumpyEncoder, self).default(obj)

def csv_to_geojson(csv_file_path: str, output_file_path: str = None) -> Dict[str, Any]:
    """
    Convert CSV file to GeoJSON format
    
    Args:
        csv_file_path: Path to the input CSV file
        output_file_path: Path for output GeoJSON file (optional)
    
    Returns:
        Dictionary containing GeoJSON data
    """
    
    try:
        # Read CSV file
        df = pd.read_csv(csv_file_path, sep=',')  # Assuming comma-separated values
        # Initialize GeoJSON structure
        geojson = {
            "type": "FeatureCollection",
            "crs": {
                "type": "name",
                "properties": {
                    "name": "urn:ogc:def:crs:OGC:1.3:CRS84"
                }
            },
            "features": []
        }
        print("CSV columns found:", df.columns.tolist())
        # Convert each row to a GeoJSON feature
        for i in range(len(df)):
            row = df.iloc[i]
            
            # Create feature with properties
            feature = {
                "type": "Feature",
                "properties": {
                    "name": str(row['name']),  # Convert to string
                    "district": str(row['district']),
                    "city": "新北市",
                    "description": str(row['description']),
                    "category": str(row['分類類別']),
                },
                "geometry": {
                    "type": "Point",
                    "coordinates": [
                        float(row['altitude']),  # Convert to regular Python float
                        float(row['latitude'])
                    ]
                }
            }
            
            geojson["features"].append(feature)
        
        # Save to file if output path is provided
        if output_file_path:
            with open(output_file_path, 'w', encoding='utf-8') as f:
                # Use custom encoder to handle NumPy types
                json.dump(geojson, f, ensure_ascii=False, indent=2, cls=NumpyEncoder)
            print(f"GeoJSON saved to: {output_file_path}")
        
        return geojson
        
    except FileNotFoundError:
        print(f"Error: CSV file '{csv_file_path}' not found")
        return None
    except Exception as e:
        print(f"Error processing CSV: {str(e)}")
        import traceback
        traceback.print_exc()  # Print full traceback for better debugging
        return None

def main():
    """
    Main function to run the conversion
    """
    input_csv = "shopping_area_ntpei.csv"
    output_geojson = "shopping_area_newtaipei.geojson"
    
    # Convert CSV to GeoJSON
    result = csv_to_geojson(input_csv, output_geojson)
    
    if result:
        print(f"Successfully converted {len(result['features'])} features")
        print("Preview of first feature:")
        if result['features']:
            print(json.dumps(result['features'][0], ensure_ascii=False, indent=2))
    else:
        print("Conversion failed")

if __name__ == "__main__":
    main()