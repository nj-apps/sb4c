package microClustering

import(
  "encoding/json"
)

type clustererJSON struct {
	//Paramètres
	McRadius float64  `json:"mc_radius"` // rayon du cluster
	MinSize  int     `json:"min_size"`// nombre minimum d'éléments composant un micro cluster pour qu'il soit pris en compte pour la génération du jeu de données
	Zones    int     `json:"zones"`// nombre de zones concentriques pour la répartition
	// statistiques sur les µC
	MediumSize       float64 `json:"medium_size"`
	SigmaSize        float64 `json:"sigma_size"`
	OutlierThreshold float64 `json:"outlier_threshold"`// un µC est considéré comme outlier si Weight < mediumSize-outlierThreshold*sigmaSize

	//structures du cluster
	VectorSize int `json:"vector_size"`
	Distance   string   `json:"distance_function"` // fonction utilisée pour évaluer les distances
	Mc         []microcluster `json:"mc_list"`// liste de tous les microclusters créés
}

type classifierJSON struct {
	Classes       map[int]clustererJSON `json:"classes"`//map avec un clusterer par classe. Le label de la classe est obligatoirement un int
	Radius        float64 `json:"radius"`
	Threshold     int   `json:"threshold"`  //Seuls les µC dont la taille dépasse le seuil de prise en compte seront utilisés pour la génération du jeu de données
	Outlier       float64 `json:"outlier"`//"nombre de sigmas, un cluster est considéré comme outlier si weight < µ(weights)-outliers*stddev(weights)"
	LabelID       int `json:"label_id"`
	Verbose       int`json:"verbose"`
	Zones         int `json:"zones"`
	CheckOutliers bool `json:"check_outliers"`// TODO: si un point est un outlier pour l'ensemble des classe alors renvoie une classe -1
}


// Export Clusterer to Json
func (c Clusterer)ToJson() ([]byte,error) {
  toExport:=clustererJSON{
McRadius: c.mcRadius,
MinSize:c.minSize,
Zones:c.zones,
MediumSize:c.mediumSize,
SigmaSize:c.sigmaSize,
OutlierThreshold:c.outlierThreshold,
VectorSize:c.vectorSize,
Distance: c.distFunction,
  }


  toExport.Mc=[]microcluster{}
  for _,v:=range c.mc {
    toExport.Mc=append(toExport.Mc,*v)
  }
  return json.Marshal(toExport)
}


// Export Clusterer to Json
func NewClustererFromJson(data []byte) (*Clusterer, error) {

  toImport:=clustererJSON{}
  err:=json.Unmarshal(data, &toImport)
if err!=nil{
  return nil, err
}

newClusterer:=Clusterer{
  mcRadius:toImport.McRadius,
  minSize:toImport.MinSize ,
  zones:toImport.Zones ,
  mediumSize:toImport.MediumSize ,
  sigmaSize:toImport.SigmaSize ,
  outlierThreshold:toImport.OutlierThreshold ,
  vectorSize:toImport.VectorSize ,
  distFunction :toImport.Distance,
}

  newClusterer.mc=[]*microcluster{}
  for _,v:=range toImport.Mc {

    mc:=microcluster{
      Weight:v.Weight,
    Zones: v.Zones,
    }
    mc.Center= make([]float64,len(v.Center))
    copy(mc.Center,v.Center)
    newClusterer.mc=append(newClusterer.mc,&mc)
  }
  return &newClusterer,nil
}

// NewClassifierFromJson creates a new classifier from to Json export
func NewClassifierFromJson(data []byte) (*Classifier, error) {

  toImport:=classifierJSON{}
  err:=json.Unmarshal(data, &toImport)
if err!=nil{
  return nil, err
}

newClassifier:=Classifier{
  CheckOutliers: toImport.CheckOutliers,
  labelID: toImport.LabelID,
  outlier:toImport.Outlier,
  Radius: toImport.Radius,
  threshold:toImport.Threshold,
  Verbose:toImport.Verbose,
  zones: toImport.Zones,
}

newClassifier.classes=make(map[int]*Clusterer)

  for k,v:=range toImport.Classes {
    d,err:=json.Marshal(v)
    if err!=nil {
      return nil, err
    }

    newClusterer,err:=NewClustererFromJson(d)
    if err!=nil {
      return nil, err
    }
    newClassifier.classes[k]=newClusterer

  }
  return &newClassifier,nil
}


func (c Clusterer)toJsonStruct() clustererJSON {
  toExport:=clustererJSON{
McRadius: c.mcRadius,
MinSize:c.minSize,
Zones:c.zones,
MediumSize:c.mediumSize,
SigmaSize:c.sigmaSize,
OutlierThreshold:c.outlierThreshold,
VectorSize:c.vectorSize,
Distance: c.distFunction,
  }


  toExport.Mc=[]microcluster{}
  for _,v:=range c.mc {
    toExport.Mc=append(toExport.Mc,*v)
  }
  return toExport
}

// Export Classifier to Json
func (c Classifier)ToJson() ([]byte,error) {
  toExport:=classifierJSON{
    CheckOutliers:c.CheckOutliers,
    Radius:c.Radius,
    Verbose:c.Verbose,
    LabelID:c.labelID,
    Outlier:c.outlier,
    Threshold:c.threshold,
    Zones:c.zones,
  }

  toExport.Classes=make(map[int]clustererJSON)
  for k,cl:=range c.classes {
    toExport.Classes[k]=(*cl).toJsonStruct()
  }

  return json.Marshal(toExport)
}
