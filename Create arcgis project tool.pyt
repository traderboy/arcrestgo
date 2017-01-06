import subprocess
import arcpy
from arcpy import env
import sqlite3
import xml.etree.ElementTree
import os
import json
import zipfile
from arcpy import mapping
import os
from xml.dom.minidom import parse
from datetime import datetime
import time
import copy
import shutil
import types
import ConfigParser
Config = ConfigParser.ConfigParser()


arcpy.env.overwriteOutput = True
#notes:  urlKey in portals.self.json must be blank or it will try to authenticate at arcgis.com
#other gotchas
#For polygon styles, makes sure to use "style": "esriSFSSolid" and NOT "style": "esriSLSSolid" for the outline style
                   


#import time
#env.workspace = "CURRENT"
#env.addOutputsToMap = False
#env.overwriteOutput = True
arcpy.env.overwriteOutput = True

toolkitPath = os.path.abspath(os.path.dirname(__file__)).replace("\\","/")

class Toolbox(object):
    def __init__(self):
        self.label =  "Create ArcServices toolbox"
        self.alias  = "arcservices"
        self.canRunInBackground = False
        # List of tool classes associated with this toolbox
        self.tools = [CreateNewProject]

class CreateNewProject(object):
    def __init__(self):
        self.label       = "Convert map document to JSON"
        self.alias="arcservices"
        self.description = "Creates the JSON files for a standalone ArcGIS Online/Server node application.  Note:  you need to fill out the project information in the File->Map Document Properties before running."
    def getParameterInfo(self):

        Config.read(toolkitPath+"/settings.ini")
        servername = arcpy.Parameter(
            displayName="Enter server FQDN (example: www.esri.com)",
            name="servername",
            datatype="GPString",
            parameterType="Required",
            direction="Input",
            multiValue=False)
        try:
           servername.value = Config.get("settings","server")
        except Exception as e:
           pass
        
        if not servername.value:
            servername.value = "gis.biz.tm"

        username = arcpy.Parameter(
            displayName="Enter your username",
            name="username",
            datatype="GPString",
            parameterType="Optional",
            direction="Input",
            multiValue=False)
        try:
           username.value= Config.get("settings","username")
        except Exception as e:
           pass

        if not username.value:
            username.value="shale"

        #projecttitle = arcpy.Parameter(
        #    displayName="Enter your project title",
        #    name="projectname",
        #    datatype="GPString",
        #    parameterType="Required",
        #    direction="Input",
        #    multiValue=False)

        #projectname = arcpy.Parameter(
        #    displayName="Enter your project name (no spaces)",
        #    name="projectname",
        #    datatype="GPString",
        #    parameterType="Required",
        #    direction="Input",
        #    multiValue=False)

        #tags = arcpy.Parameter(
        #    displayName="Enter tags",
        #    name="tags",
        #    datatype="GPString",
        #    parameterType="Optional",
        #    direction="Input",
        #    multiValue=False)
        #
        #summary = arcpy.Parameter(
        #    displayName="Enter project summary",
        #    name="summary",
        #    datatype="GPString",
        #    parameterType="Optional",
        #    direction="Input",
        #    multiValue=False)
        #
        #description = arcpy.Parameter(
        #    displayName="Enter project description",
        #    name="description",
        #    datatype="GPString",
        #    parameterType="Optional",
        #    direction="Input",
        #    multiValue=False)

        outputfolder = arcpy.Parameter(
            displayName="Enter output folder",
            name="outputfolder",
            datatype="DEFolder",
            parameterType="Required",
            direction="Input")
        try:
            outputfolder.value= Config.get("settings","destination")
        except Exception as e:
           pass
        
        if not outputfolder.value:
            outputfolder.value="C:/massappraisal/agolnode/hpl12"

        #param0.filter.type = "ValueList"
        #param0.filter.list = ["Street","Aerial","Terrain","Topographic"]
        parameters = [servername,username,outputfolder]
        #username,projecttitle,projectname,tags,summary,description,
        return parameters

    def isLicensed(self): #optional
        return True

    def updateParameters(self, parameters): #optional
       if parameters[2].altered:
          try:
              os.makedirs(parameters[2].valueAsText)
          except Exception as e:
              return
       return

    def updateMessages(self, parameters): #optional
        return

    def execute(self, parameters, messages):
        serverName = parameters[0].valueAsText
        username = parameters[1].valueAsText        
        baseDestinationPath = parameters[2].valueAsText
        created_ts=int(time.time()*1000)

        # suppose you want to add it to the current MXD (open MXD)
        try:
           if type(messages)==types.StringType:
              vals = messages.split("|")
              if len(vals)>1:
                 serverName = vals[1]
              if len(vals)>2:
                 username= vals[2]
              if len(vals)>3:
                 baseDestinationPath=vals[3].replace("\\","/")
              mxdName=vals[0].replace("\\","/")
              mxd = arcpy.mapping.MapDocument(mxdName)
           else:
              mxd = arcpy.mapping.MapDocument("CURRENT")
        except Exception as e:
           printMessage("Still Unable to open map document.  Make sure background processing is unchecked in the geoprocessing options")
           return

        printMessage("Exporting dataframe: " + mxd.activeDataFrame.name)
        serviceName = mxd.activeDataFrame.name.replace(" ","").lower()
        if serviceName=='Layers':
           printMessage("Rename the dataframe from Layers to service name.  Must be valid service name (no spaces)")
           return

        #mxd.makeThumbnail ()
        #toolkitPath = os.path.abspath(os.path.dirname(__file__)).replace("\\","/")
        
        templatePath = toolkitPath + "/templates"
        if not os.path.exists(templatePath):
            printMessage("Template path not found: " + templatePath)
            return

        
        cfgfile = open(toolkitPath+"/settings.ini",'w')
        try:
            Config.add_section("settings")
        except Exception as e:
            pass

        printMessage("Server name: " +serverName)
        printMessage("User name: " + username)
        printMessage("Destination path: " + baseDestinationPath)

        Config.set("settings","server",serverName)
        Config.set("settings","username",username)
        Config.set("settings","destination",baseDestinationPath)
        Config.write(cfgfile)
        cfgfile.close()
        del cfgfile       
        
        if baseDestinationPath:
              baseDestinationPath = unicode(baseDestinationPath).encode('unicode-escape')
              baseDestinationPath=baseDestinationPath.replace("\\","/")
        else:
              baseDestinationPath = toolkitPath

        #baseDestinationPath = baseDestinationPath + os.sep + serviceName
        serviceDestinationPath = baseDestinationPath + os.sep + serviceName

        #if the folder does not exist create it
        if not os.path.exists(baseDestinationPath):
            os.makedirs(serviceDestinationPath)
        else:
            #check to see if service already exists.  If so, remove it so it can be overwritten
            if os.path.exists(serviceDestinationPath):
               try:
                  shutil.rmtree(serviceDestinationPath)
               except Exception as e:
                  printMessage("Unable to remove destination path")
                  return
            try:
               os.makedirs(serviceDestinationPath)
            except Exception as e:
               printMessage("Unable to create destination path")


        servicesDestinationPath = serviceDestinationPath + "/services"
        if not os.path.exists(servicesDestinationPath):
            os.makedirs(servicesDestinationPath)
        printMessage("Services path: " +servicesDestinationPath)

        dataDestinationPath = serviceDestinationPath + "/shapefiles"
        if not os.path.exists(dataDestinationPath):
            os.makedirs(dataDestinationPath)
        printMessage("Shapefile path: " +dataDestinationPath)

        replicaDestinationPath = serviceDestinationPath + "/replicas"
        if not os.path.exists(replicaDestinationPath):
            os.makedirs(replicaDestinationPath)
        printMessage("Replica path: " +replicaDestinationPath)

        mapfileDestinationPath = serviceDestinationPath + "/mapfiles"
        if not os.path.exists(mapfileDestinationPath):
            os.makedirs(mapfileDestinationPath)
        printMessage("Mapfile path: " +mapfileDestinationPath)

        symbols = getSymbology(mxd)
        dataFrames = arcpy.mapping.ListDataFrames(mxd, "*")

        if os.path.exists(baseDestinationPath + "/config.json"):
            config=openJSON(baseDestinationPath + "/config.json")
            try:
               config["services"][serviceName]={}
            except:
               printMessage("Service already exists: " + serviceName)
            config["services"][serviceName]["layers"]={}
        else:
           config={}
           config["services"]={}
           config["services"][serviceName]={"layers":{}}

        config["hostname"]=serverName
        config["username"]=username
        config["mxd"]=mxd.filePath
        

        #config["services"][serviceName]["layers"]={}

        fullname = mxd.author
        if fullname=="":
           printMessage("Author missing in File->Map Document Properties")
           return
        first_name = fullname.split(' ')[0]
        last_name = fullname.split(' ')[1]
        email_address = first_name + '.' + last_name + '@' + serverName
        if not username:
           username=fullname.lower().replace(" ","")

        title  = mxd.title
        if title=="":
           printMessage("Title missing in File->Map Document Properties")
           return

        tags  = mxd.tags
        if not tags:
           tags=""

        summary  = mxd.summary
        if not summary:
           summary=""

        description  = mxd.description
        if not description:
           description=""

        if not os.path.exists(baseDestinationPath + "/portals.self.json"):
           portals_self_json=openJSON(templatePath + "/portals.self.json")
           portals_self_json['portalHostname']=serverName
           portals_self_json['defaultExtent']['xmin']=mxd.activeDataFrame.extent.XMin
           portals_self_json['defaultExtent']['ymin']=mxd.activeDataFrame.extent.YMin
           portals_self_json['defaultExtent']['xmax']=mxd.activeDataFrame.extent.XMax
           portals_self_json['defaultExtent']['ymax']=mxd.activeDataFrame.extent.YMax
           portals_self_json['user']['fullName']=fullname
           portals_self_json['user']['firstName']=first_name
           portals_self_json['user']['lastName']=last_name
           portals_self_json['user']['email']=email_address
           portals_self_json['user']['username']=username
           saveJSON(baseDestinationPath + "/portals.self.json",portals_self_json)

        if not os.path.exists(baseDestinationPath + "/community.users.json"):
           community_users_json=openJSON(templatePath + "/community.users.json")
           community_users_json['fullName']=fullname
           community_users_json['firstName']=first_name
           community_users_json['lastName']=last_name
           community_users_json['email']=email_address
           community_users_json['username']=username
           community_users_json['created']=created_ts
           community_users_json['modified']=created_ts
           community_users_json['lastLogin']=created_ts
           #community_users_json['groups'][0]['userMembership']['username']=username
           saveJSON(baseDestinationPath + "/community.users.json",community_users_json)

        #User info
        content_users_json=openJSON(templatePath + "/content.users.json")
        content_users_json['username']=username
        #content_users_json['items'][0]['created']=int(time.time()*1000)
        saveJSON(baseDestinationPath + "/content.users.json",content_users_json)

        #Search results
        if not os.path.exists(baseDestinationPath + "/search.json"):
            search_json=openJSON(templatePath + "/search.json")
            #search_json['results'][0]=username
            baseResult = search_json['results'][0]
            search_json['results']=[]
        else:
            search_json=openJSON(baseDestinationPath + "/search.json")
            baseResult = search_json['results'][0]
            #see if result already exists and delete it
            for idx, val in enumerate(search_json['results']):
               if val["id"] == serviceName:
                  del search_json['results'][idx]

            #search_json['results']

        #add stuff for each dataframe below

        #community groups
        #community_groups_json=openJSON(templatePath + "/community.groups.json")
        #saveJSON(destinationPath + "/community.groups.json",community_groups_json)
        shutil.copy2(templatePath + "/community.groups.json", baseDestinationPath + "/community.groups.json")
        #os.system("copy "+ templatePath + "/community.groups.json " + servicesDestinationPath + "/community.groups.json")
        #result = 0

        feature_services={"currentVersion":10.3,"folders":[],"services":[]}
        #if not os.path.exists(servicesDestinationPath+"/FeatureServer.json"):
        #     saveJSON(servicesDestinationPath + "/FeatureServer.json",response)
        #else:
        #     featureServer_json=openJSON(servicesDestinationPath + "/FeatureServer.json")
        #     if not serviceName in featureServer_json['folders']:
        #         featureServer_json['folders'].append(serviceName);
        #         saveJSON(servicesDestinationPath + "/FeatureServer.json",featureServer_json)
        #     #create base FeatureServer.json file with folders for each service
        #     #,"folders":["Canvas","Demographics","Elevation","Ocean","Polar","Reference","Specialty","Utilities"]

        #for dataFrame in dataFrames:
        if mxd.activeDataFrame:
           dataFrame = mxd.activeDataFrame
           serviceName = dataFrame.name.replace(" ","").lower()
           #mxd.activeDataFrame.name
           if serviceName=='Layers':
              printMessage("Rename the dataframe from Layers to service name.  Must be valid service name (no spaces)")
              return

           #else:
           #   dataFrame = dataFrame #mxd.activeDataFrame
           operationalLayers = []
           operationalTables = []
           operationalTablesObj = []
           allData=[]

           #for df in arcpy.mapping.ListDataFrames(mxd):
           for lyr in arcpy.mapping.ListLayers(mxd, "", dataFrame):
              # Exit if the current layer is not a service layer.
              if lyr.isServiceLayer or lyr.supports("SERVICEPROPERTIES"):  # or not lyr.visible
                continue
              #lyr.visible=True
              #opLayer = {
              #    "id": lyr.name,
              #    "title": lyr.name,
              #    "url": lyr.serviceProperties["Resturl"]+ "/" + lyr.longName + "/" + lyr.serviceProperties["ServiceType"],
              #    "opacity": (100 - lyr.transparency) / 100,
              #    "visibility": lyr.visible
              #}
              printMessage("Exporting layer: " + lyr.name)
              operationalLayers.append(lyr)
              allData.append(lyr)
              #arcpy.mapping.RemoveLayer(df, lyr)

           if len(operationalLayers)==0:
              printMessage("No Feature layers found in data frame!")
              return

           id=len(operationalLayers)
           for tbl in arcpy.mapping.ListTableViews(mxd, "", dataFrame):
              operationalTables.append(tbl)
              allData.append(tbl)
              operationalTablesObj.append({"name":tbl.name,"id":id})
              id=id+1

           #lyrpath=os.getcwd().replace("\\","/")
           #lyrpath = os.path.abspath(os.path.dirname(__file__)).replace("\\","/")

           ext = operationalLayers[0].getExtent()
           dataFrame.extent = ext

           desc = arcpy.Describe(operationalLayers[0])
           if hasattr(desc, "layer"):
                 ws=desc.layer.catalogPath.replace("\\","/")
           else:
                 ws=os.path.dirname(desc.catalogPath).replace("\\","/")

           for j,rel in enumerate(allData):
              printMessage(str(j) + ": " + rel.name)

           relationships = [c.name for c in arcpy.Describe(ws).children if c.datatype == "RelationshipClass"]

           relArr=[]
           desc = arcpy.Describe(lyr)
           if not desc.relationshipClassNames:
              return rel
           if hasattr(desc, "layer"):
                 featureName=os.path.basename(desc.layer.catalogPath)
                 rootFGDB=desc.layer.catalogPath.replace("\\","/")
           else:
                 featureName=os.path.basename(desc.catalogPath)
                 rootFGDB=os.path.dirname(desc.catalogPath).replace("\\","/")

           config["fgdb"]=rootFGDB
           relationshipObj = {}
           relations={}
           #for index in xrange(0, field_info.count):
           #[u'farm_tracts_inspections__ATTACHREL', u'farm_tractsInspectionRelClass', u'homesites_inspections__ATTACHREL', u'homesitesInspectionRelClass', u'grazing_inspections__ATTACHREL', u'grazing_permitteeRelClass', u'grazing_permitteesInspectionRelClass']
           #for j,rel in enumerate(desc.relationshipClassNames):
           id=0
           printMessage("Find relationships")
           for rc in relationships:
             relDesc = arcpy.Describe(rootFGDB+"/"+rc)
             if relDesc.isAttachmentRelationship:
                  continue
             #printMessage("Relationship Name: " + rc)
             #printMessage("Origin Class Names")
             #printMessage(relDesc.originClassNames)

             #printMessage("Origin Class Keys")
             #printMessage(relDesc.originClassKeys)

             #printMessage("Destination Class Names")
             #printMessage(relDesc.destinationClassNames)

             #printMessage("Destination Class Keys")
             #printMessage(relDesc.destinationClassKeys)

             #printMessage("Key type:  "+relDesc.keyType)
             #printMessage(relDesc.notification)
             #printMessage("backwardPathLabel:  "+relDesc.backwardPathLabel)
             #printMessage("forwardPathLabel:  "+relDesc.forwardPathLabel)

             originId=getDataIndex(allData,relDesc.originClassNames[0])
             destId=getDataIndex(allData,relDesc.destinationClassNames[0])

             relatedTableId=0
             role=""
             key=""
             relations[str(id)]={"oTable":relDesc.originClassNames[0],"dTable":relDesc.destinationClassNames[0],"oJoinKey":relDesc.originClassKeys[0][0],"dJoinKey":relDesc.originClassKeys[1][0],"oId":originId,"dId":destId}

             relObj = {"id":id,"name":relDesc.forwardPathLabel,"relatedTableId":destId,"cardinality":"esriRelCardinality"+relDesc.cardinality,"role":"esriRelRoleOrigin","keyField":relDesc.originClassKeys[0][0],"composite":relDesc.isComposite}
             id=id+1

             try:
                len(relationshipObj[relDesc.originClassNames[0]])
             except:
                relationshipObj[relDesc.originClassNames[0]]=[]

             relationshipObj[relDesc.originClassNames[0]].append(relObj)

             try:
                len(relationshipObj[relDesc.destinationClassNames[0]])
             except:
                relationshipObj[relDesc.destinationClassNames[0]]=[]

             relObj = {"id":id,"name":relDesc.backwardPathLabel,"relatedTableId":originId,"cardinality":"esriRelCardinality"+relDesc.cardinality,"role":"esriRelRoleDestination","keyField":relDesc.originClassKeys[1][0],"composite":relDesc.isComposite}
             relationshipObj[relDesc.destinationClassNames[0]].append(relObj)

           #printMessage(json.dumps(relationshipObj, indent=4, sort_keys=True))
           config["services"][serviceName]["relationships"]=relations
           #return



           #printMessage(relationships)
           #for rc in relationships:
           #   rc_path = ws + "\\" + rc
           #   des_rc = arcpy.Describe(rc_path)
           #   printMessage(des_rc.originClassNames)

           #rc_list = [c.name for c in arcpy.Describe(workspace).children if c.datatype == "RelationshipClass"]
           #for rc in rc_list:
           #rc_path = workspace + "\\" + rc
           #des_rc = arcpy.Describe(rc_path)
           #origin = des_rc.originClassNames
           #destination = des_rc.destinationClassNames

           #mxd.activeDataFrame=dataFrame
           mxd.activeView = dataFrame.name
           arcpy.RefreshActiveView()

           #out_file_name = r"c:\thumbnails\{basename}.png".format(basename=os.path.basename(featureclass))
           # Export "thumbnail" of data frame

           #if the folder does not exist create it
           if not os.path.exists(servicesDestinationPath+"/thumbnails/"):
              os.makedirs(servicesDestinationPath+"/thumbnails/")

           out_file_name = servicesDestinationPath + "/thumbnails/" + serviceName + ".png"
           arcpy.mapping.ExportToPNG(mxd, out_file_name, dataFrame, 200, 133)

           #dataFrame = arcpy.mapping.ListDataFrames(mxd, "*")[0]
           #if dataFrame != mxd.activeDataFrame:
           #   printMessage("Active data frame is not the first data frame")

           feature_services['folders'].append(serviceName);

           #now set path to serviceName folder
           #destinationPath = servicesDestinationPath + "/data" #+ serviceName
           #print destinationPath
           #printMessage("Spatial JSON destination path: " + servicesDestinationPath)
           #if the folder does not exist create it
           #if not os.path.exists(destinationPath):
           #    os.makedirs(destinationPath)

           rootService_json={"folders": [], "services":[{"name":serviceName,"type":"FeatureServer","url":"http://"+serverName+"/rest/services/"+serviceName+"/FeatureServer"},{"name":serviceName,"type":"MapServer"}], "currentVersion": 10.3}
           saveJSON(servicesDestinationPath + "/"+serviceName+".json",rootService_json)

           #analysis = arcpy.mapping.AnalyzeForMSD(mxd)
           #
           #for key in ('messages', 'warnings', 'errors'):
           #  printMessage( "----" + key.upper() + "---")
           #  vars = analysis[key]
           #  for ((message, code), layerlist) in vars.iteritems():
           #    printMessage( "    " + message + " (CODE %i)" % code)
           #    printMessage( "       applies to:")
           #    for layer in layerlist:
           #        printMessage( layer.name)
           #    printMessage("")

#           sddraft = templatePath + serviceName + '.sddraft'
#           sd = templatePath + serviceName + '.sd'
#           summary = 'Sample output'
#           tags = 'county, counties, population, density, census'
#
#           # create service definition draft
#           analysis = arcpy.mapping.CreateMapSDDraft(mxd, sddraft, serviceName, 'ARCGIS_SERVER')
#
#           for key in ('messages', 'warnings', 'errors'):
#             printMessage("----" + key.upper() + "---")
#             vars = analysis[key]
#             for ((message, code), layerlist) in vars.iteritems():
#               printMessage("    " +  message + " (CODE %i)" % code)
#               printMessage("       applies to:")
#               for layer in layerlist:
#                   printMessage(layer.name)
#               printMessage("")
#
#           printMessage("")
#           printMessage("")
#           #arcpy.StageService_server(sddraft, sd)
#
#           # stage and upload the service if the sddraft analysis did not contain errors
#           if analysis['errors'] == {}:
#               # Execute StageService
#               arcpy.StageService_server(sddraft, sd)
#               # Execute UploadServiceDefinition
#               #arcpy.UploadServiceDefinition_server(sd, con)
#           else:
#               # if the sddraft analysis contained errors, display them
#               #arcpy.StageService_server(sddraft, sd)
#               printMessage(analysis['errors'])
#               #print analysis['errors']

           #arcpy.mapping.ConvertToMSD(mxd,toolkitPath+"/output.msd",dataFrame, "NORMAL", "NORMAL")
           #mxde = MxdExtras(mxd)
           #for lyr in mxde.itervalues():
           #   printMessage("Layer Name: " + lyr.name )
           #   printMessage("Layer Symbology Field Name: " + lyr.symbologyFieldName)

           oldspatialref = dataFrame.spatialReference
           coordinateSystem = 'GEOGCS["GCS_WGS_1984",DATUM["D_WGS_1984",SPHEROID["WGS_1984",6378137.0,298.257223563]],PRIMEM["Greenwich",0.0],UNIT["Degree",0.017453292519943295]]'
           #set to wgs84
           dataFrame.spatialReference = coordinateSystem
           #get coors of extent center in new coordinate system
           x = (dataFrame.extent.XMin + dataFrame.extent.XMax)/2
           y = (dataFrame.extent.YMin + dataFrame.extent.YMax)/2
           #printMessage(str(dataFrame.extent.XMin) + "," + str(dataFrame.extent.YMin) + "," + str(dataFrame.extent.XMax)  + "," + str(dataFrame.extent.YMax))
           xmin_geo=dataFrame.extent.XMin
           xmax_geo=dataFrame.extent.XMax
           ymin_geo=dataFrame.extent.YMin
           ymax_geo=dataFrame.extent.YMax
           # set dataframe spatial ref back
           dataFrame.spatialReference = oldspatialref

           output = {
             "extent": {
               "xmin": dataFrame.extent.XMin,
               "ymin": dataFrame.extent.YMin,
               "xmax": dataFrame.extent.XMax,
               "ymax": dataFrame.extent.YMax
             },
             "scale": dataFrame.scale,
             "rotation": dataFrame.rotation,
             "spatialReference": {"wkid": dataFrame.spatialReference.PCSCode}
           }

           result=copy.deepcopy(baseResult) # deep  copy
           result['snippet']=summary
           result['title']=dataFrame.description
           result['id']=serviceName
           #result['extent']=[0,0]
           result['extent'][0]=[0,0]
           result['extent'][1]=[0,0]
           result['extent'][0][0]=xmin_geo
           result['extent'][0][1]=ymin_geo
           result['extent'][1][0]=xmax_geo
           result['extent'][1][1]=ymax_geo
           result['owner']=username
           result['created']=created_ts
           result['modified']=created_ts
           if tags!="":
               result['tags']=tags.split(",")
           search_json['results'].append(result)
           #result = result + 1

           #only need to update the operationalLayers
           content_items_json=openJSON(templatePath + "/content.items.data.json")
           opLayers = content_items_json['operationalLayers']=getOperationalLayers(operationalLayers,serverName,serviceName)
           opTables = content_items_json['tables']=getTables(operationalTables,serverName,serviceName,len(opLayers))

           saveJSON(servicesDestinationPath + "/content.items.data.json",content_items_json)

           content_items_json=openJSON(templatePath + "/content.items.json")
           content_items_json["id"]=title
           content_items_json["owner"]=username
           content_items_json["created"]=created_ts
           content_items_json["modified"]=created_ts
           content_items_json["title"]=title
           content_items_json["snippet"]=summary
           content_items_json["description"]=description
           content_items_json['extent'][0][0]=xmin_geo
           content_items_json['extent'][0][1]=ymin_geo
           content_items_json['extent'][1][0]=xmax_geo
           content_items_json['extent'][1][1]=ymax_geo

           content_items_json["type"]="Feature Service"
           content_items_json["url"]="http://"+serverName+"/rest/services/"+serviceName+"/FeatureServer"
           saveJSON(servicesDestinationPath + "/content.items.json",content_items_json)

           #create JSON description of all services.  Each dataframe is a service for this application.
           featureserver_json={
              "currentVersion":10.3,
              "services": [{
                 "name":serviceName,
                 "type":"FeatureServer",
                 "url": "http://"+serverName + "/arcgis/rest/services/"+serviceName+"/FeatureServer"
              }]
           }
           saveJSON(servicesDestinationPath + "/FeatureServer.json",featureserver_json)

           #create JSON description of all layers in the service.
           featureserver_json=openJSON(templatePath + "/name.FeatureServer.json")
           featureserver_json['initialExtent']['xmin']=dataFrame.extent.XMin
           featureserver_json['initialExtent']['ymin']=dataFrame.extent.YMin
           featureserver_json['initialExtent']['xmax']=dataFrame.extent.XMax
           featureserver_json['initialExtent']['ymax']=dataFrame.extent.YMax
           featureserver_json['fullExtent']['xmin']=dataFrame.extent.XMin
           featureserver_json['fullExtent']['ymin']=dataFrame.extent.YMin
           featureserver_json['fullExtent']['xmax']=dataFrame.extent.XMax
           featureserver_json['fullExtent']['ymax']=dataFrame.extent.YMax
           featureserver_json['layers'] = getLayers(operationalLayers)
           featureserver_json['tables']=operationalTablesObj
           saveJSON(servicesDestinationPath + "/FeatureServer.json",featureserver_json)

           maps_json=openJSON(templatePath + "/name.MapServer.json")
           maps_json['initialExtent']['xmin']=dataFrame.extent.XMin
           maps_json['initialExtent']['ymin']=dataFrame.extent.YMin
           maps_json['initialExtent']['xmax']=dataFrame.extent.XMax
           maps_json['initialExtent']['ymax']=dataFrame.extent.YMax
           maps_json['fullExtent']['xmin']=dataFrame.extent.XMin
           maps_json['fullExtent']['ymin']=dataFrame.extent.YMin
           maps_json['fullExtent']['xmax']=dataFrame.extent.XMax
           maps_json['fullExtent']['ymax']=dataFrame.extent.YMax
           maps_json['layers'] = featureserver_json['layers']
           maps_json['server']=serverName
           maps_json['name']=serviceName
           maps_json['mapName']=serviceName
           maps_json['tables']=operationalTablesObj
           saveJSON(servicesDestinationPath + "/MapServer.json",maps_json)

           minx=str(dataFrame.extent.XMin)
           miny=str(dataFrame.extent.YMin)
           maxx=str(dataFrame.extent.XMax)
           maxy=str(dataFrame.extent.YMax)
           createReplica(mxd,dataFrame,allData,replicaDestinationPath,toolkitPath,username,serviceName,serverName,minx,miny,maxx,maxy)

           #create a JSON service file for each feature layer -- broken ---
           serviceRep=[]

           id=0
           for lyr in operationalLayers:
               desc = arcpy.Describe(lyr)
               if hasattr(desc, "layer"):
                   featureName=os.path.basename(desc.layer.catalogPath)
               else:
                   featureName=os.path.basename(desc.catalogPath)
               printMessage(lyr.name+": " + featureName)

               feature_json=openJSON(templatePath + "/name.FeatureServer.id.json")
               feature_json['defaultVisibility']=lyr.visible
               feature_json['description'] = lyr.description
               feature_json['fields']=getFields(lyr)
               #type=esriFieldTypeOID

               #for i in feature_json:
               #     printMessage(i + ": " + str(feature_json[i]))
               #printMessage(feature_json['displayField'])
               #if lyr.showLabels:
               lbl=""
               if lyr.supports("LABELCLASSES"):
                   for lblclass in lyr.labelClasses:
                       lblclass.showClassLabels = True
                       #feature_json.displayField
                       lbl=lblclass.expression.replace("[","").replace("]","")
               #lblclass.expression = " [Label]"
               if lbl!="":
                  feature_json['displayField']=lbl
               else:
                  feature_json['displayField']=getDisplayField(feature_json['fields'])

               if desc.shapeType:
                   if desc.shapeType=='Polygon':
                      feature_json['geometryType']='esriGeometryPolygon'
                      feature_json['templates'][0]['drawingTool']="esriFeatureEditToolPolygon"
                   elif desc.shapeType=='Polyline':
                      feature_json['geometryType']='esriGeometryPolyline'
                      feature_json['templates'][0]['drawingTool']="esriFeatureEditToolPolyline"
                   elif desc.shapeType=='Point':
                      feature_json['geometryType']='esriGeometryPoint'

                   elif desc.shapeType=='MultiPoint':
                      feature_json['geometryType']='esriGeometryMultiPoint'
               feature_json['id']=id
               feature_json['name']=lyr.name
               if desc.hasOID:
                   feature_json['objectIdField']=desc.OIDFieldName
                   feature_json['objectIdFieldName']=desc.OIDFieldName
               if desc.hasGlobalID:
                   feature_json['globalIdField'] = desc.globalIDFieldName
               else:
                   feature_json['globalIdField'] = ""
               feature_json['indexes']=getIndexes(lyr)
               feature_json['minScale']=lyr.minScale
               feature_json['maxScale']=lyr.maxScale
               #bad below, should be Feature Layer, not FeatureLayer
               #feature_json['type']=desc.dataType #'Feature Layer'
               feature_json['extent']['xmin']=desc.extent.XMin
               feature_json['extent']['ymin']=desc.extent.YMin
               feature_json['extent']['xmax']=desc.extent.XMax
               feature_json['extent']['ymax']=desc.extent.YMax
               #feature_json['indexes']=[]
               feature_json['templates'][0]['name']=serviceName
               attributes={}
               for field in feature_json['fields']:
                   #printMessage(field['name'])
                   if field['editable']:
                      attributes[ field['name'] ]=None
               feature_json['templates'][0]['prototype']['attributes']=attributes

               #feature_json['drawingInfo']['renderer']['symbol']=getSymbol(lyr)
               #feature_json['relationships']=getRelationships(lyr,id,len(operationalLayers),operationalTables,relationshipObj)
               feature_json['relationships']=relationshipObj[featureName] #getRelationships(lyr,relationshipObj)

               feature_json['drawingInfo']=getSymbol(lyr,symbols[featureName],lyr.name)
               #set editor tracking fields
               editorTracking={}
               if desc.editorTrackingEnabled:
                  editorTracking['creationDateField']=desc.createdAtFieldName
                  editorTracking['creatorField']=desc.creatorFieldName
                  editorTracking['editDateField']=desc.editedAtFieldName
                  editorTracking['editorField']=desc.editorFieldName
                  feature_json['editFieldsInfo']=editorTracking
               else:
                  del feature_json['editFieldsInfo']

               feature_json['editingInfo']={"lastEditDate":created_ts}

               if arcpy.Exists(rootFGDB+"/"+featureName+"__ATTACH"):
                  feature_json['hasAttachments']=True
               else:
                  feature_json['hasAttachments']=False

               #getSymbol(lyr,symbols[featureName],lyr.name)
               #opLayers = content_items_json['operationalLayers']=getOperationalLayers(operationalLayers,serverName,serviceName)
               saveJSON(servicesDestinationPath + "/FeatureServer."+str(id)+".json",feature_json)

               #now create a MapServer json file
               mapserver_json=openJSON(templatePath + "/name.MapServer.id.json")
               mapserver_json['indexes']=feature_json['indexes']
               mapserver_json['extent']=feature_json['extent']
               mapserver_json['fields']=feature_json['fields']
               
               mapserver_json['templates']=feature_json['templates']
               mapserver_json['drawingInfo']=feature_json['drawingInfo']
               mapserver_json['geometryType']=feature_json['geometryType']

               saveJSON(servicesDestinationPath + "/MapServer."+str(id)+".json",feature_json)

               #save replica file
               feature_json=openJSON(templatePath + "/name.FeatureServer.id.json")

               #steps: save layer to blank mxd, save it, run arcpy.CreateRuntimeContent on mxd
               createSingleReplica(templatePath,dataFrame,lyr,replicaDestinationPath,toolkitPath,feature_json,serverName,serviceName,username,id)
               #save mapserver .map file
               saveMapfile(mapfileDestinationPath + "/"+lyr.name+".map",lyr,desc,dataDestinationPath,mapserver_json)

               id = id+1

           #create a JSON geometry file for each feature layer
           id=0
           for lyr in operationalLayers:
               desc = arcpy.Describe(lyr)

               fSet = arcpy.FeatureSet()
               fSet.load(desc.dataElement.catalogPath)
               fdesc = arcpy.Describe(fSet)
               #printMessage(fdesc.json)
               dataName = os.path.basename(desc.dataElement.catalogPath)
               layerObj={"name":lyr.name,"data":dataName}

               feature_json = json.loads(fdesc.json)
               #feature_json=openJSON(templatePath + "/name.FeatureServer.id.query.json")
               #feature_json['features']=getFeatures(lyr)
               #feature_json['fields']=getFields(lyr)
               #if desc.shapeType:
               #    if desc.shapeType=='Polygon':
               #       feature_json['geometryType']='esriGeometryPolygon'
               #    elif desc.shapeType=='Polyline':
               #       feature_json['geometryType']='esriGeometryPolyline'
               #    elif desc.shapeType=='Point':
               #       feature_json['geometryType']='esriGeometryPoint'
               #    elif desc.shapeType=='MultiPoint':
               #       feature_json['geometryType']='esriGeometryMultiPoint'

               if desc.hasOID:
                   feature_json['objectIdField']=desc.OIDFieldName
                   layerObj["oidname"]=desc.OIDFieldName
                   feature_json['objectIdFieldName']=desc.OIDFieldName
               if desc.hasGlobalID:
                   feature_json['globalIdField'] = desc.globalIDFieldName
                   layerObj["globaloidname"]=desc.globalIDFieldName

               layerObj["type"]="layer"
               printMessage("Saving layer " + str(id) + " to JSON")
               saveJSON(servicesDestinationPath + "/FeatureServer."+str(id)+".query.json",feature_json)

               #create a JSON OBJECTID file used in ArcGIS for showing the attribute table
               #remove all fields except OBJECTID
               #feature_json['fields']=[{"alias":"OBJECTID","name":"OBJECTID","type":"esriFieldTypeInteger","alias":"OBJECTID","sqlType":"sqlTypeOther","defaultValue":None,"domain":None}]
               feature_json['fields']=[{"alias":"OBJECTID","name":"OBJECTID","type":"esriFieldTypeOID","sqlType":"sqlTypeOther","defaultValue":None,"domain":None,"nullable":False,"editable":False}]
               features=[]
               #for i in feature_json['fields']:
               #   if i['name'] != 'OBJECTID':
               #      del i
               #      #del feature_json['fields'][i]
               for i in feature_json['features']:
                  features.append({"attributes":{"OBJECTID":i['attributes']['OBJECTID']}})
               feature_json['features']=features
                  #for j in i['attributes']:
                  #    if j == 'OBJECTID':
                  #      attribute={"OBJECTID":j}
                  #      #del j
                  #for j in feature_json['features'][i]['attributes']:
                  #   if feature_json['features'][i]['attributes'][j]['name'] != 'OBJECTID':
                  #      del feature_json.features[i]['attributes'][j]

               saveJSON(servicesDestinationPath + "/FeatureServer."+str(id)+".objectid.json",feature_json)
               layerObj["itemId"]= lyr.name.replace(" ","_")+str(id)
               config["services"][serviceName]["layers"][str(id)]=layerObj
               id = id+1

           #now save any tables
           for tbl in operationalTables:
               desc = arcpy.Describe(tbl)
               featureName=os.path.basename(desc.catalogPath)
               feature_json=openJSON(templatePath + "/name.RecordSet.id.json")
               #feature_json['description'] = tbl.description
               tableObj={"name":tbl.name,"data":featureName}
               feature_json['fields']=getFields(tbl)
               feature_json['displayField']=getDisplayField(feature_json['fields'])

               #feature_json['relationships']=getRelationships(tbl,id,len(operationalLayers),operationalTables,relationshipObj)
               #feature_json['relationships']=getRelationships(lyr,relationshipObj)

               feature_json['id']=id
               feature_json['name']=tbl.name
               if desc.hasOID:
                   feature_json['objectIdField']=desc.OIDFieldName
                   feature_json['objectIdFieldName']=desc.OIDFieldName
                   tableObj["oidname"]=desc.OIDFieldName
               if desc.hasGlobalID:
                   feature_json['globalIdField'] = desc.globalIDFieldName
                   tableObj["globaloidname"]=desc.globalIDFieldName
               else:
                   feature_json['globalIdField'] = ""
               tableObj["type"]="table"
               feature_json['indexes']=getIndexes(tbl)
               feature_json['templates'][0]['name']=serviceName
               attributes={}
               for field in feature_json['fields']:
                   #printMessage(field['name'])
                   if field['editable']:
                      attributes[ field['name'] ]=None
               feature_json['templates'][0]['prototype']['attributes']=attributes

               printMessage(tbl.name+": " + featureName)
               feature_json['relationships']=relationshipObj[featureName]

               #set editor tracking fields
               editorTracking={}
               if desc.editorTrackingEnabled:
                  editorTracking['creationDateField']=desc.createdAtFieldName
                  editorTracking['creatorField']=desc.creatorFieldName
                  editorTracking['editDateField']=desc.editedAtFieldName
                  editorTracking['editorField']=desc.editorFieldName
                  feature_json['editFieldsInfo']=editorTracking
               else:
                  del feature_json['editFieldsInfo']

               feature_json['editingInfo']={"lastEditDate":created_ts}

               if arcpy.Exists(rootFGDB+"/"+featureName+"__ATTACH"):
                  feature_json['hasAttachments']=True
               else:
                  feature_json['hasAttachments']=False

               saveJSON(servicesDestinationPath + "/FeatureServer."+str(id)+".json",feature_json)
               tableObj["itemId"]= tbl.name.replace(" ","_")+str(id)
               config["services"][serviceName]["layers"][str(id)]=tableObj

               fSet = arcpy.RecordSet()
               fSet.load(desc.catalogPath)
               fdesc = arcpy.Describe(fSet)
               #printMessage(fdesc.json)
               feature_json = json.loads(fdesc.json)
               #dataName = os.path.basename(desc.dataElement.catalogPath)
               #layerObj={"name":lyr.name,"data":dataName}
               printMessage("Saving table " + str(id) + " to JSON")
               saveJSON(servicesDestinationPath + "/FeatureServer."+str(id)+".query.json",feature_json)

               id = id+1

           #export all layers to shapefiles for rendering in mapserver
           for lyr in operationalLayers:
               desc = arcpy.Describe(lyr)
               if desc.dataType == "FeatureLayer":
                  printMessage("Saving layer to shapefile: "+ lyr.name)
                  arcpy.FeatureClassToFeatureClass_conversion(desc.dataElement.catalogPath,
                     dataDestinationPath,
                     lyr.name+".shp")

           id = id+1

        #now save the search results
        search_json['total']=len(search_json['results'])
        saveJSON(baseDestinationPath + "/search.json",search_json)
        #save root FeatureServer.json file
        saveJSON(baseDestinationPath + "/FeatureServer.json",feature_services)
        saveJSON(baseDestinationPath + "/MapServer.json",feature_services)
        saveJSON(baseDestinationPath + "/config.json",config)
        printMessage("Finished")



def openJSON(name):
    printMessage("Loading JSON template: " +name)
    with open(name, "r+") as f:
       json_data = json.load(f)
       f.close()
       return json_data

def saveJSON(name,json_data): #optional
    with open(name,'w') as f:
       f.write(json.dumps(json_data))

def clearSelections(mxd):
    for df in arcpy.mapping.ListDataFrames(mxd):
       for lyr in arcpy.mapping.ListLayers(mxd, "", df):
          # Exit if the current layer is not a service layer.
          if  lyr.isServiceLayer or lyr.supports("SERVICEPROPERTIES"):  # or not lyr.visible
            continue
          if not lyr.isFeatureLayer:
            continue
          printMessage(lyr.name +": " + arcpy.Describe(lyr).catalogPath)
          arcpy.SelectLayerByAttribute_management(lyr, "CLEAR_SELECTION")
          #arcpy.Describe(lyr).catalogPath

def getSymbology(mxd):
    msdPath = os.path.abspath(os.path.dirname(__file__)).replace("\\","/")+"/output.msd"
    #msdPath = self.mxdPath.replace(self.MXD_SUFFIX, self.MSD_SUFFIX)

    # Delete temporary msd if it exists
    if os.path.exists(msdPath):
        os.remove(msdPath)

    clearSelections(mxd)
    arcpy.mapping.ConvertToMSD(mxd,msdPath)
    symbols={}

    zz = zipfile.ZipFile(msdPath)
    EXCLUDED_FILE_NAMES = ["DocumentInfo.xml", "GISProject.xml", "layers/layers.xml"]
    for fileName in (fileName for fileName in zz.namelist() if not fileName in EXCLUDED_FILE_NAMES):
        printMessage("Opening: " + fileName)
        dom = parse(zz.open(fileName))
        symb = dom.getElementsByTagName("Symbolizer")
        if symb.length>0:
            name=fileName.split(".")[0]
            rootname = name.split("/")
            if len(rootname)>1:
              name=rootname[1]
            printMessage("Symbology found for : " + name + " length: " + str(symb.length))
            symbols[name]=symb
            #printMessage("Found: " + str(symb.length))
            #name,lyr = self.loadMsdLayerDom(dom)
            #if name != "":
            #   self[name] = lyr
    del zz
    return symbols

def getLayers(opLayers):
  layers=[]
  id=0
  for lyr in opLayers:
     layer={
        "name":lyr.name,
        "id":id,
        "parentLayerId":-1,
        "defaultVisibility":lyr.visible,
        "subLayerIds":None,
        "minScale":lyr.minScale,
        "maxScale":lyr.maxScale
      }
     #{"name":"farm_tracts","type":"FeatureServer","url":"http://services5.arcgis.com/KOH6W4WICH5gzytg/ArcGIS/rest/services/farm_tracts/FeatureServer"},
     layers.append(layer)
     id=id+1
  return layers

def createSingleReplica(templatePath,df,lyr,replicaDestinationPath,toolkitPath,feature_json,serverName,serviceName,username,id):
  blankmxd = arcpy.mapping.MapDocument(templatePath + "/blank.mxd")
  df = arcpy.mapping.ListDataFrames(blankmxd)[0]
  arcpy.mapping.AddLayer(df, lyr)

  tmpMxd=templatePath+"/temp.mxd"
  if os.path.exists(tmpMxd):
     os.remove(tmpMxd)

  blankmxd.saveACopy(tmpMxd)
  #mxd.save()
  desc=arcpy.Describe(lyr)

  saveReplica(tmpMxd,replicaDestinationPath + "/"+lyr.name,lyr,desc)
  #move to root folder and delete the "data" folder
  filenames = next(os.walk(replicaDestinationPath + "/"+lyr.name+"/data/"))[2]
  printMessage("Rename " + replicaDestinationPath + "/"+lyr.name+"/data/"+filenames[0]+" to "+ replicaDestinationPath+"/"+lyr.name+".geodatabase")
  #if offline geodatabase exists, must delete first
  newReplicaDB=replicaDestinationPath+"/"+lyr.name+".geodatabase"
  try:
     if os.path.exists(newReplicaDB):
        os.rmdir(newReplicaDB)
  except:
     printMessage("Unable to remove old replica geodatabase")

  os.rename(replicaDestinationPath + "/"+lyr.name+"/data/"+filenames[0], newReplicaDB)
  try:
     os.rmdir(replicaDestinationPath + "/"+lyr.name+"/data/")
     os.rmdir(replicaDestinationPath + "/"+lyr.name)
  except:
     printMessage("Unable to remove replica folders")
  #if os.path.exists(tmpMxd):
  #   os.remove(tmpMxd)
  #sqliteDb=replicaDestinationPath + "/"+lyr.name+"/data/"+serviceName+".geodatabase"
  #sqliteDb=replicaDestinationPath + "/"+lyr.name+".geodatabase"

  serviceRep=[]
  if os.path.exists(newReplicaDB):
     #dom = parse(templatePath+"/replica.xml")
     #xml = createXML(dom,serverName,serviceName,lyr.name):
     ret=updateReplicaPaths(newReplicaDB,lyr.name,feature_json,"http://"+serverName + "/arcgis/rest/services/"+serviceName+"/FeatureServer",serverName,serviceName,username,id)
     ret1 = updateReplicaPaths(newReplicaDB,lyr.name,feature_json,"http://"+serverName + "/arcgis/rest/services/"+serviceName+"/FeatureServer",serverName,serviceName,username,id)

     for i in ret1:
        serviceRep.append(i)

     name=replicaDestinationPath + "/"+lyr.name+".sql"
     with open(name,'w') as f:
          f.write("SELECT load_extension( 'stgeometry_sqlite.dll', 'SDE_SQL_funcs_init');\n")
          for i in ret:
             f.write(i)
             f.write(";\n")

          #f.write(ret[0])
          #f.write(";\n")
          #f.write(ret[1])
          #f.write(";\n")
          #f.write(ret[2])
          #f.write(";\n")
          #f.write(ret[3])
          #f.write(";\n")

          f.close()
     printMessage("Running " + toolkitPath+"/spatialite/spatialite.exe " + newReplicaDB + "  < " + name)
     try:
        os.system(toolkitPath+"/spatialite/spatialite.exe " + newReplicaDB + "  < " + name)
     except:
        printMessage("Unable to run sql commands")
     #C:\tmp\32bit>spatialite \\192.168.2.124\d\git\arcnodegis\root\replica\rangeland_units.geodatabase < \\192.168.2.124\d\git\arcnodegis\root\replica\rangeland_units.sql

def saveReplica(tmpMxd,replicaPath,lyr,desc):
   arcpy.CreateRuntimeContent_management(tmpMxd,
      replicaPath,
      lyr.name,"#","#",
      "FEATURE_AND_TABULAR_DATA","NON_OPTIMIZE_SIZE","ONLINE","PNG","1","#")
   printMessage("Saved replica: "+ replicaPath)

def updateReplicaPaths(replicaPath,layerName,feature_json,servicePath,serverName,serviceName,username,id):
   minx=str(feature_json['extent']['xmin'])
   miny=str(feature_json['extent']['ymin'])
   maxx=str(feature_json['extent']['xmax'])
   maxy=str(feature_json['extent']['ymax'])

   sql1=('INSERT INTO GDB_Items("ObjectID", "UUID", "Type", "Name", "PhysicalName", "Path", "Url", "Properties", "Defaults", "DatasetSubtype1", "DatasetSubtype2", "DatasetInfo1", "DatasetInfo2", "Definition", "Documentation", "ItemInfo", "Shape")'
    " select MAX(ObjectID)+1, '{7B6EB064-7BF6-42C8-A116-2E89CD24A000}', '{5B966567-FB87-4DDE-938B-B4B37423539D}', 'MyReplica', 'MYREPLICA', 'MyReplica', '', 1, NULL, NULL, NULL, "
    "'http://"+serverName+"/arcgis/rest/services/"+serviceName+"/FeatureServer', '"+username+"',"
    "'<GPSyncReplica xsi:type=''typens:GPSyncReplica'' xmlns:xsi=''http://www.w3.org/2001/XMLSchema-instance'' xmlns:xs=''http://www.w3.org/2001/XMLSchema'' xmlns:typens=''http://www.esri.com/schemas/ArcGIS/10.3''>"
    "<ReplicaName>MyReplica</ReplicaName><ID>"+str(id)+"</ID><ReplicaID>{7b6eb064-7bf6-42c8-a116-2e89cd24a000}</ReplicaID>"
    "<ServiceName>http://"+serverName+"/arcgis/rest/services/"+serviceName+"/FeatureServer</ServiceName>"
    "<Owner>"+username+"</Owner>"
    "<Role>esriReplicaRoleChild</Role><SyncModel>esriSyncModelPerLayer</SyncModel><Direction>esriSyncDirectionBidirectional</Direction><CreationDate>2015-09-02T13:48:33</CreationDate><LastSyncDate>1970-01-01T00:00:01</LastSyncDate>"
    "<ReturnsAttachments>true</ReturnsAttachments><SpatialRelation>esriSpatialRelIntersects</SpatialRelation><QueryGeometry xsi:type=''typens:PolygonN''><HasID>false</HasID><HasZ>false</HasZ><HasM>false</HasM><Extent xsi:type=''typens:EnvelopeN''>"
    "<XMin>"+minx+"</XMin><YMin>"+miny+"</YMin><XMax>"+maxx+"</XMax><YMax>"+maxy+"</YMax></Extent><RingArray xsi:type=''typens:ArrayOfRing''><Ring xsi:type=''typens:Ring''>"
    "<PointArray xsi:type=''typens:ArrayOfPoint''>"
    "<Point xsi:type=''typens:PointN''><X>"+minx+"</X><Y>"+miny+"</Y></Point><Point xsi:type=''typens:PointN''><X>"+maxx+"</X><Y>"+miny+"</Y></Point>"
    "<Point xsi:type=''typens:PointN''><X>"+maxx+"</X><Y>"+maxy+"</Y></Point><Point xsi:type=''typens:PointN''><X>"+minx+"</X><Y>"+maxy+"</Y></Point>"
    "<Point xsi:type=''typens:PointN''><X>"+minx+"</X><Y>"+miny+"</Y></Point></PointArray></Ring></RingArray>"
    "<SpatialReference xsi:type=''typens:ProjectedCoordinateSystem''><WKT>PROJCS[&quot;WGS_1984_Web_Mercator_Auxiliary_Sphere&quot;,GEOGCS[&quot;GCS_WGS_1984&quot;,DATUM[&quot;D_WGS_1984&quot;,SPHEROID[&quot;WGS_1984&quot;,6378137.0,298.257223563]],PRIMEM[&quot;Greenwich&quot;,0.0],UNIT[&quot;Degree&quot;,0.0174532925199433]],PROJECTION[&quot;Mercator_Auxiliary_Sphere&quot;],PARAMETER[&quot;False_Easting&quot;,0.0],PARAMETER[&quot;False_Northing&quot;,0.0],PARAMETER[&quot;Central_Meridian&quot;,0.0],PARAMETER[&quot;Standard_Parallel_1&quot;,0.0],PARAMETER[&quot;Auxiliary_Sphere_Type&quot;,0.0],UNIT[&quot;Meter&quot;,1.0],AUTHORITY[&quot;EPSG&quot;,3857]]</WKT><XOrigin>-20037700</XOrigin><YOrigin>-30241100</YOrigin><XYScale>10000</XYScale><ZOrigin>-100000</ZOrigin><ZScale>10000</ZScale><MOrigin>-100000</MOrigin><MScale>10000</MScale><XYTolerance>0.001</XYTolerance><ZTolerance>0.001</ZTolerance><MTolerance>0.001</MTolerance><HighPrecision>true</HighPrecision><WKID>102100</WKID><LatestWKID>3857</LatestWKID></SpatialReference></QueryGeometry>"

    "<GPSyncDatasets xsi:type=''typens:ArrayOfGPSyncDataset''>"
    "<GPSyncDataset xsi:type=''typens:GPSyncDataset''><DatasetID>5</DatasetID><DatasetName>"+layerName+"</DatasetName><DatasetType>esriDTFeatureClass</DatasetType>"
    "<LayerID>"+str(id)+"</LayerID><LayerName>"+layerName+"</LayerName><Direction>esriSyncDirectionBidirectional</Direction><ReplicaServerGen xsi:type=''xs:long''>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type=''xs:long''>1000</ReplicaClientDownloadGen>"
    "<ReplicaClientUploadGen xsi:type=''xs:long''>1000</ReplicaClientUploadGen><ReplicaClientAcknowledgeUploadGen xsi:type=''xs:long''>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>true</IncludeRelated>"
    "<QueryOption>esriRowsTypeFilter</QueryOption><HasAttachments>true</HasAttachments></GPSyncDataset>"
    "<GPSyncDataset xsi:type=''typens:GPSyncDataset''><DatasetID>6</DatasetID><DatasetName>"+layerName+"__ATTACH</DatasetName><DatasetType>esriDTTable</DatasetType>"
    "<LayerID>0</LayerID>"
    "<LayerName>"+layerName+"</LayerName><Direction>esriSyncDirectionBidirectional</Direction>"
    "<ReplicaServerGen xsi:type=''xs:long''>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type=''xs:long''>1000</ReplicaClientDownloadGen><ReplicaClientUploadGen xsi:type=''xs:long''>1000</ReplicaClientUploadGen>"
    "<ReplicaClientAcknowledgeUploadGen xsi:type=''xs:long''>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>false</IncludeRelated><QueryOption>esriRowsTypeFilter</QueryOption><IsAttachment>true</IsAttachment></GPSyncDataset>"
    "</GPSyncDatasets></GPSyncReplica>'"
    ", NULL, NULL, NULL from GDB_Items")

   sql2=('INSERT INTO GDB_Items("ObjectID", "UUID", "Type", "Name", "PhysicalName", "Path", "Url", "Properties", "Defaults", "DatasetSubtype1", "DatasetSubtype2", "DatasetInfo1", "DatasetInfo2", "Definition", "Documentation", "ItemInfo", "Shape")'
    " select MAX(ObjectID)+1, '{AE8D3C7E-9890-4BF4-B946-5BE50A1CC279}', '{D86502F9-9758-45C6-9D23-6DD1A0107B47}', '"+layerName+"', '"+layerName.upper()+"', 'MyReplica"+str(id)+"\\"+layerName+"', '', 1, NULL, NULL, NULL, NULL, NULL, "
    "'<GPSyncDataset xsi:type=''typens:GPSyncDataset'' xmlns:xsi=''http://www.w3.org/2001/XMLSchema-instance'' xmlns:xs=''http://www.w3.org/2001/XMLSchema'' xmlns:typens=''http://www.esri.com/schemas/ArcGIS/10.3''>"
    "<DatasetID>5</DatasetID>"
    "<DatasetName>"+layerName+"</DatasetName>"
    "<DatasetType>esriDTFeatureClass</DatasetType><LayerID>0</LayerID><LayerName>"+layerName+"</LayerName><Direction>esriSyncDirectionBidirectional</Direction>"
    "<ReplicaServerGen xsi:type=''xs:long''>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type=''xs:long''>1000</ReplicaClientDownloadGen><ReplicaClientUploadGen xsi:type=''xs:long''>1000</ReplicaClientUploadGen>"
    "<ReplicaClientAcknowledgeUploadGen xsi:type=''xs:long''>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>true</IncludeRelated><QueryOption>esriRowsTypeFilter</QueryOption>"
    "<HasAttachments>true</HasAttachments></GPSyncDataset>', NULL, NULL, NULL from GDB_Items")

   sql3=('INSERT INTO GDB_Items("ObjectID", "UUID", "Type", "Name", "PhysicalName", "Path", "Url", "Properties", "Defaults", "DatasetSubtype1", "DatasetSubtype2", "DatasetInfo1", "DatasetInfo2", "Definition", "Documentation", "ItemInfo", "Shape")'
   " select MAX(ObjectID)+1, '{55C5E7E4-834D-4D44-A12C-991E7F8B4645}', '{D86502F9-9758-45C6-9D23-6DD1A0107B47}', '"+layerName+"__ATTACH', '"+layerName.upper()+"__ATTACH', 'MyReplica_"+str(id)+"\\"+layerName+"__ATTACH', '', 1, NULL, NULL, NULL, NULL, NULL, "
   "'<GPSyncDataset xsi:type=''typens:GPSyncDataset'' xmlns:xsi=''http://www.w3.org/2001/XMLSchema-instance'' xmlns:xs=''http://www.w3.org/2001/XMLSchema'' xmlns:typens=''http://www.esri.com/schemas/ArcGIS/10.3''>"
   "<DatasetID>6</DatasetID>"
   "<DatasetName>"+layerName+"__ATTACH</DatasetName><DatasetType>esriDTTable</DatasetType><LayerID>0</LayerID><LayerName>"+layerName+"</LayerName><Direction>esriSyncDirectionBidirectional</Direction>"
   "<ReplicaServerGen xsi:type=''xs:long''>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type=''xs:long''>1000</ReplicaClientDownloadGen><ReplicaClientUploadGen xsi:type=''xs:long''>1000</ReplicaClientUploadGen>"
   "<ReplicaClientAcknowledgeUploadGen xsi:type=''xs:long''>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>false</IncludeRelated><QueryOption>esriRowsTypeFilter</QueryOption>"
   "<IsAttachment>true</IsAttachment></GPSyncDataset>',"
   " NULL, NULL, NULL from GDB_Items")

   #sql4='update "GDB_ServiceItems" set "ItemInfo"=replace("ItemInfo",|"capabilities":"Query"|,|"capabilities":"Create,Delete,Query,Update,Editing,Sync"|)'
   #sql4=sql4.replace("|","'")
   #sql5='update "GDB_Items" set ObjectId=ROWID'

   return [sql1,sql2,sql3]

   #sqliteReplicaPaths(sql1,sql2,sql3)

def sqliteReplicaPaths(sql1,sql2,sql3):
   conn = sqlite3.connect(replicaPath)
   conn.enable_load_extension(True)
   c = conn.cursor()
   sql4='update "GDB_ServiceItems" set "ItemInfo"=replace("ItemInfo",|"capabilities":"Query"|,|"capabilities":"Create,Delete,Query,Update,Editing,Sync"|)'
   sql4=sql4.replace("|","'")

   # A) Inserts an ID with a specific value in a second column
   #http://services5.arcgis.com/KOH6W4WICH5gzytg/ArcGIS/rest/services/rangeland_units/FeatureServer
   try:
      #c.execute("update GDB_Items set DatasetInfo1='"+servicePath + "',DatasetInfo2='"+username+"',Definition='+xml.toxml()+' where Name='"+name+"'")
      printMessage("SELECT load_extension( 'c:\Program Files (x86)\ArcGIS\Desktop10.3\DatabaseSupport\SQLite\Windows32\stgeometry_sqlite.dll', 'SDE_SQL_funcs_init')")
      c.execute("SELECT load_extension( 'c:/Program Files (x86)/ArcGIS/Desktop10.3/DatabaseSupport/SQLite/Windows32/stgeometry_sqlite.dll', 'SDE_SQL_funcs_init')")
      #c.execute("SELECT load_extension( 'c:/tmp/32bit/stgeometry_sqlite.dll', 'SDE_SQL_funcs_init')")
      #c.execute("SELECT load_extension( 'c:/Program Files (x86)/ArcGIS/Desktop10.3/DatabaseSupport/SQLite/Windows64/stgeometry_sqlite.dll', 'SDE_SQL_funcs_init')")
   except sqlite3.IntegrityError:
      printMessage("Error in sql integrity")
   else:
      printMessage("Error in sql")

   try:
      #c.execute("update GDB_Items set DatasetInfo1='"+servicePath + "',DatasetInfo2='"+username+"',Definition='+xml.toxml()+' where Name='"+name+"'")
      printMessage(sql1)
      c.execute(sql1)
   except sqlite3.IntegrityError:
      printMessage("Error in sql integrity")
   else:
      printMessage("Error in sql")

   try:
      printMessage(sql2)
      c.execute(sql2)
   except sqlite3.IntegrityError:
      printMessage("Error in sql integrity")
   else:
      printMessage("Error in sql")

   try:
      printMessage(sql3)
      c.execute(sql3)
   except sqlite3.IntegrityError:
      printMessage("Error in sql integrity")
   else:
      printMessage("Error in sql")

   try:
      printMessage(sql4)
      c.execute(sql4)
   except sqlite3.IntegrityError:
      printMessage("Error in sql integrity")
   else:
      printMessage("Error in sql")

   conn.commit()
   conn.close()

def __createXML(xmlFile,serverName,serviceName,layerName):
  dom = parse(zz.open(xmlFile))
  symb = dom.getElementsByTagName("Symbolizer")
  dom.getElementsByTagName("ServiceName")
  dom.getElementsByTagName("ReplicaName")
  dom.getElementsByTagName("ID")
  dom.getElementsByTagName("ReplicaID")
  dom.getElementsByTagName("Owner")
  key.firstChild.data = "new text"

#DatasetID
#DatasetName
#DatasetType>esriDTFeatureClass</DatasetType>
#LayerID
#LayerName

  # Open original file
  #et = xml.etree.ElementTree.parse(xmlFile)

  # Append new tag: <a x='1' y='abc'>body text</a>
  #new_tag = xml.etree.ElementTree.SubElement(et.getroot(), 'a')
  #new_tag.text = 'body text'
  #new_tag.attrib['x'] = '1' # must be str; cannot be an int
  #new_tag.attrib['y'] = 'abc'

  # Write back to file
  #et.write('file.xml')
  #et.write('file_new.xml')


def __updateReplica(layer):
  #8	8	{7B6EB064-7BF6-42C8-A116-2E89CD24A000}	{5B966567-FB87-4DDE-938B-B4B37423539D}	MyReplica_7	MYREPLICA_7	MyReplica_7		1	NULL	NULL	NULL	http://services5.arcgis.com/KOH6W4WICH5gzytg/ArcGIS/rest/services/rangeland_units/FeatureServer	hplgis	<GPSyncReplica xsi:type='typens:GPSyncReplica' xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance' xmlns:xs='http://www.w3.org/2001/XMLSchema' xmlns:typens='http://www.esri.com/schemas/ArcGIS/10.3'><ReplicaName>MyReplica_7</ReplicaName><ID>7</ID><ReplicaID>{7b6eb064-7bf6-42c8-a116-2e89cd24a000}</ReplicaID><ServiceName>http://services5.arcgis.com/KOH6W4WICH5gzytg/ArcGIS/rest/services/rangeland_units/FeatureServer</ServiceName><Owner>hplgis</Owner><Role>esriReplicaRoleChild</Role><SyncModel>esriSyncModelPerLayer</SyncModel><Direction>esriSyncDirectionBidirectional</Direction><CreationDate>2015-09-02T13:48:33</CreationDate><LastSyncDate>1970-01-01T00:00:01</LastSyncDate><ReturnsAttachments>true</ReturnsAttachments><SpatialRelation>esriSpatialRelIntersects</SpatialRelation><QueryGeometry xsi:type='typens:PolygonN'><HasID>false</HasID><HasZ>false</HasZ><HasM>false</HasM><Extent xsi:type='typens:EnvelopeN'><XMin>"+minx+"</XMin><YMin>4210044.6521089608</YMin><XMax>-12229227.023144344</XMax><YMax>4369440.8825313309</YMax></Extent><RingArray xsi:type='typens:ArrayOfRing'><Ring xsi:type='typens:Ring'><PointArray xsi:type='typens:ArrayOfPoint'><Point xsi:type='typens:PointN'><X>"+minx+"</X><Y>4210044.6521089608</Y></Point><Point xsi:type='typens:PointN'><X>-12229227.023144344</X><Y>4210044.6521089608</Y></Point><Point xsi:type='typens:PointN'><X>-12229227.023144344</X><Y>4369440.8825313309</Y></Point><Point xsi:type='typens:PointN'><X>"+minx+"</X><Y>4369440.8825313309</Y></Point><Point xsi:type='typens:PointN'><X>"+minx+"</X><Y>4210044.6521089608</Y></Point></PointArray></Ring></RingArray><SpatialReference xsi:type='typens:ProjectedCoordinateSystem'><WKT>PROJCS[&quot;WGS_1984_Web_Mercator_Auxiliary_Sphere&quot;,GEOGCS[&quot;GCS_WGS_1984&quot;,DATUM[&quot;D_WGS_1984&quot;,SPHEROID[&quot;WGS_1984&quot;,6378137.0,298.257223563]],PRIMEM[&quot;Greenwich&quot;,0.0],UNIT[&quot;Degree&quot;,0.0174532925199433]],PROJECTION[&quot;Mercator_Auxiliary_Sphere&quot;],PARAMETER[&quot;False_Easting&quot;,0.0],PARAMETER[&quot;False_Northing&quot;,0.0],PARAMETER[&quot;Central_Meridian&quot;,0.0],PARAMETER[&quot;Standard_Parallel_1&quot;,0.0],PARAMETER[&quot;Auxiliary_Sphere_Type&quot;,0.0],UNIT[&quot;Meter&quot;,1.0],AUTHORITY[&quot;EPSG&quot;,3857]]</WKT><XOrigin>-20037700</XOrigin><YOrigin>-30241100</YOrigin><XYScale>10000</XYScale><ZOrigin>-100000</ZOrigin><ZScale>10000</ZScale><MOrigin>-100000</MOrigin><MScale>10000</MScale><XYTolerance>0.001</XYTolerance><ZTolerance>0.001</ZTolerance><MTolerance>0.001</MTolerance><HighPrecision>true</HighPrecision><WKID>102100</WKID><LatestWKID>3857</LatestWKID></SpatialReference></QueryGeometry><GPSyncDatasets xsi:type='typens:ArrayOfGPSyncDataset'><GPSyncDataset xsi:type='typens:GPSyncDataset'><DatasetID>5</DatasetID><DatasetName>rangeland_units</DatasetName><DatasetType>esriDTFeatureClass</DatasetType><LayerID>0</LayerID><LayerName>rangeland_units</LayerName><Direction>esriSyncDirectionBidirectional</Direction><ReplicaServerGen xsi:type='xs:long'>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type='xs:long'>1000</ReplicaClientDownloadGen><ReplicaClientUploadGen xsi:type='xs:long'>1000</ReplicaClientUploadGen><ReplicaClientAcknowledgeUploadGen xsi:type='xs:long'>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>true</IncludeRelated><QueryOption>esriRowsTypeFilter</QueryOption><HasAttachments>true</HasAttachments></GPSyncDataset><GPSyncDataset xsi:type='typens:GPSyncDataset'><DatasetID>6</DatasetID><DatasetName>rangeland_units__ATTACH</DatasetName><DatasetType>esriDTTable</DatasetType><LayerID>0</LayerID><LayerName>rangeland_units</LayerName><Direction>esriSyncDirectionBidirectional</Direction><ReplicaServerGen xsi:type='xs:long'>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type='xs:long'>1000</ReplicaClientDownloadGen><ReplicaClientUploadGen xsi:type='xs:long'>1000</ReplicaClientUploadGen><ReplicaClientAcknowledgeUploadGen xsi:type='xs:long'>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>false</IncludeRelated><QueryOption>esriRowsTypeFilter</QueryOption><IsAttachment>true</IsAttachment></GPSyncDataset></GPSyncDatasets></GPSyncReplica>	NULL	NULL	NULL
  #9	9	{AE8D3C7E-9890-4BF4-B946-5BE50A1CC279}	{D86502F9-9758-45C6-9D23-6DD1A0107B47}	rangeland_units	RANGELAND_UNITS	MyReplica_7\rangeland_units		1	NULL	NULL	NULL	NULL	NULL	<GPSyncDataset xsi:type='typens:GPSyncDataset' xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance' xmlns:xs='http://www.w3.org/2001/XMLSchema' xmlns:typens='http://www.esri.com/schemas/ArcGIS/10.3'><DatasetID>5</DatasetID><DatasetName>rangeland_units</DatasetName><DatasetType>esriDTFeatureClass</DatasetType><LayerID>0</LayerID><LayerName>rangeland_units</LayerName><Direction>esriSyncDirectionBidirectional</Direction><ReplicaServerGen xsi:type='xs:long'>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type='xs:long'>1000</ReplicaClientDownloadGen><ReplicaClientUploadGen xsi:type='xs:long'>1000</ReplicaClientUploadGen><ReplicaClientAcknowledgeUploadGen xsi:type='xs:long'>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>true</IncludeRelated><QueryOption>esriRowsTypeFilter</QueryOption><HasAttachments>true</HasAttachments></GPSyncDataset>	NULL	NULL	NULL                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            #tables to update GDB_ColumnRegistry, GDB_ItemRelationshipTypes,GDB_Items,GDB_Layers,GDB_ServiceItems,GDB_TableRegistry,st_geometry_columns,st_spindex__rangeland_units_SHAPE,st_spindex__rangeland_units_SHAPE_rowid
  #10	10	{55C5E7E4-834D-4D44-A12C-991E7F8B4645}	{D86502F9-9758-45C6-9D23-6DD1A0107B47}	rangeland_units__ATTACH	RANGELAND_UNITS__ATTACH	MyReplica_7\rangeland_units__ATTACH		1	NULL	NULL	NULL	NULL	NULL	<GPSyncDataset xsi:type='typens:GPSyncDataset' xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance' xmlns:xs='http://www.w3.org/2001/XMLSchema' xmlns:typens='http://www.esri.com/schemas/ArcGIS/10.3'><DatasetID>6</DatasetID><DatasetName>rangeland_units__ATTACH</DatasetName><DatasetType>esriDTTable</DatasetType><LayerID>0</LayerID><LayerName>rangeland_units</LayerName><Direction>esriSyncDirectionBidirectional</Direction><ReplicaServerGen xsi:type='xs:long'>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type='xs:long'>1000</ReplicaClientDownloadGen><ReplicaClientUploadGen xsi:type='xs:long'>1000</ReplicaClientUploadGen><ReplicaClientAcknowledgeUploadGen xsi:type='xs:long'>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>false</IncludeRelated><QueryOption>esriRowsTypeFilter</QueryOption><IsAttachment>true</IsAttachment></GPSyncDataset>	NULL	NULL	NULL
  import sqlite3

  sqlite_file = 'my_first_db.sqlite'
  table_name = 'my_table_2'
  id_column = 'my_1st_column'
  column_name = 'my_2nd_column'

  # Connecting to the database file
  conn = sqlite3.connect(sqlite_file)
  c = conn.cursor()

  # A) Inserts an ID with a specific value in a second column
  try:
      c.execute("INSERT INTO {tn} ({idf}, {cn}) VALUES (123456, 'test')".\
          format(tn=table_name, idf=id_column, cn=column_name))
  except sqlite3.IntegrityError:
      printMessage('ERROR: ID already exists in PRIMARY KEY column {}'.format(id_column))

  # B) Tries to insert an ID (if it does not exist yet)
  # with a specific value in a second column
  c.execute("INSERT OR IGNORE INTO {tn} ({idf}, {cn}) VALUES (123456, 'test')".\
          format(tn=table_name, idf=id_column, cn=column_name))

  # C) Updates the newly inserted or pre-existing entry
  c.execute("UPDATE {tn} SET {cn}=('Hi World') WHERE {idf}=(123456)".\
          format(tn=table_name, cn=column_name, idf=id_column))

  conn.commit()
  conn.close()


def copyGeodatabase():
  # Execute CopyRuntimeGdbToFileGdb
  arcpy.CopyRuntimeGdbToFileGdb_conversion("D:\runtimedata\replica.geodatabase", 'replica_Copy.gdb')

def createSqliteDb():
   # Set local variables
   sqlite_database_path = "C:\data\example.gpkg"

   # Execute CreateSQLiteDatabase
   arcpy.gp.CreateSQLiteDatabase(sqlite_database_path, "GEOPACKAGE")
   # Set environment settings
   arcpy.env.workspace = "C:/data"

   # Set local variables
   outWorkspace = "c:/output/output.gdb"

   # Use ListFeatureClasses to generate a list of shapefiles in the
   #  workspace shown above.
   fcList = arcpy.ListFeatureClasses()

   # Execute CopyFeatures for each input shapefile
   for shapefile in fcList:
       # Determine the new output feature class path and name
       outFeatureClass = os.path.join(outWorkspace, shapefile.strip(".shp"))
       arcpy.CopyFeatures_management(shapefile, outFeatureClass)

   mxd = arcpy.mapping.MapDocument(r"D:\xTemp\sample.mxd")
   df = arcpy.mapping.ListDataFrames(mxd, "*")[0]
   mxd.saveACopy(r"d:\xTemp\name.mxd")
   mxd.save()

   arcpy.CreateRuntimeContent("D:/Geoprocessing/OfflineMapping/sandiego_locators_basemap.mxd",
   "D:/Geoprocessing/Results/RRuntimeContent_sandiego",
   "MyBasemapLayer","#","#",
   "FEATURE_AND_TABULAR_DATA","NON_OPTIMIZE_SIZE","ONLINE","PNG","1","#")

def saveMapfile(mapfile,lyr,desc,dataDestinationPath,mapserver_json):
  mapsize="400 400"
  symbol=""
  type=""
  style=""

  try:
     size=str(mapserver_json['drawingInfo']['renderer']['symbol']['size'])
  except:  # AttributeError:
     size="1"

  try:
     width=str(mapserver_json['drawingInfo']['renderer']['symbol']['width'])
  except: # AttributeError:
     width="1"

  try:
     color=str(mapserver_json['drawingInfo']['renderer']['symbol']['color'][0]) +" "+ str(mapserver_json['drawingInfo']['renderer']['symbol']['color'][1])+" " + str(mapserver_json['drawingInfo']['renderer']['symbol']['color'][2])
  except:  # AttributeError:
     color="0 0 0"
  try:
     outlinecolor=str(mapserver_json['drawingInfo']['renderer']['symbol']['outline']['color'][0]) +" "+ str(mapserver_json['drawingInfo']['renderer']['symbol']['outline']['color'][1])+" " + str(mapserver_json['drawingInfo']['renderer']['symbol']['outline']['color'][2])
  except:  # AttributeError:
     outlinecolor="0 0 0"

  if mapserver_json['geometryType']=='esriGeometryPolygon':
     type="POLYGON"
     style=("COLOR "+color+" \n"
       "OUTLINECOLOR "+outlinecolor+" \n"
       "WIDTH "+width+" \n")
  elif mapserver_json['geometryType']=='esriGeometryPolyline':
     type="LINE"
     style=("COLOR "+color+" \n"
       "OUTLINECOLOR "+outlinecolor+" \n"
       "WIDTH "+width+" \n")
  elif mapserver_json['geometryType']=='esriGeometryPoint':
     symbol=("SYMBOL\n"
      "NAME 'circle'\n"
      "TYPE ellipse\n"
      "FILLED true\n"
      "POINTS\n"
      "1 1\n"
      "END\n"
      "END\n")

     type="POINT"
     style=("COLOR "+color+"\n"
      "SYMBOL 'circle'\n"
      "SIZE "+size+"\n")


  mapstr = ("MAP \n "
      "NAME '" +lyr.name + "' \n"
      "STATUS ON \n"
      "EXTENT " + str(mapserver_json['extent']['xmin']) + " " + str(mapserver_json['extent']['ymin']) + " " + str(mapserver_json['extent']['xmax']) + " " + str(mapserver_json['extent']['ymax']) + "\n"
      "SIZE " + mapsize + "\n"
      "IMAGECOLOR 255 255 255 \n"
      "TRANSPARENT on \n"
      +symbol+
      "LAYER\n"
      "NAME "+lyr.name+"\n"
      "STATUS ON \n"
      "DATA 'data/"+lyr.name+"'\n"
      "TYPE "+type+"\n"
      "CLASS\n"
      "NAME '"+lyr.name+"'\n"
      "STYLE\n"
      +style+
      "END\n"
      "END\n"
      "END\n"
      "END")

  try:
    file = open(mapfile,'w')   # Trying to create a new file or open one
    file.write(mapstr)
    file.close()
  except:
    printMessage("Unable to create mapfile: " + mapstr)


def getOperationalLayers(opLayers,serverName,serviceName):
   layers=[]
   id = 0
   for lyr in opLayers:
     opLayer = {
         "id": str(id),
         "itemId": lyr.name.replace(" ","_")+str(id),
         "layerType":"ArcGISFeatureLayer",
         "title": lyr.name,
         "url": "http://"+serverName + "/arcgis/rest/services/"+serviceName+"/FeatureServer/"+str(id),
         "popupInfo": getPopupInfo(lyr),
         "opacity": (100.0 - float(lyr.transparency)) / 100.0,
         "capabilities": "Create,Delete,Query,Update,Editing,Sync",
         "visibility": lyr.visible
     }
     lbl=""
     if lyr.supports("LABELCLASSES"):
         for lblclass in lyr.labelClasses:
             lblclass.showClassLabels = True
             lbl=lblclass.expression.replace("[","").replace("]","")
     if lbl!="":
          opLayer['popupInfo']['title']=lyr.name + ":  {" + lbl + "}"
     #"opacity": (100 - lyr.transparency) / 100,
     #"url": lyr.serviceProperties["Resturl"]+ "/" + lyr.longName + "/" + lyr.serviceProperties["ServiceType"],

     id=id+1
     layers.append(opLayer)
   return layers

def getTables(opTables,serverName,serviceName,id=0):
   tbls=[]
   for tbl in opTables:
     opTable = {
         "id": str(id),
         "itemId": tbl.name.replace(" ","_")+str(id),
         #"layerType":"ArcGISFeatureLayer",
         "title": tbl.name,
         "url": "http://"+serverName + "/arcgis/rest/services/"+serviceName+"/FeatureServer/"+str(id),
         "popupInfo": getPopupInfo(tbl),
         #"visibility": lyr.visible,
         "capabilities": "Create,Delete,Query,Update,Editing,Sync"
     }
     id=id+1
     tbls.append(opTable)
   return tbls

def getDisplayField(fields):
    displayField=""
    for field in fields:
      #printMessage(field['type'])
      if field['type']=="esriFieldTypeString":
          return field['name']
    #type=esriFieldTypeOID
    return displayField

def getDisplayFieldName(lyr):
    desc = arcpy.Describe(lyr)
    if desc.dataType == "FeatureLayer":
        # Create a fieldinfo object
        fields= arcpy.ListFields(desc.dataElement.catalogPath)
    else:
        fields= arcpy.ListFields(desc.catalogPath)

    displayField=""
    for field in fields:
      #printMessage(field.type)
      if field.type=="String":
          return field.name + " {"+field.name +"}"
    #type=esriFieldTypeOID
    return displayField


def getPopupInfo(lyr):

   popInfo = {'title': getDisplayFieldName(lyr),
              'description':None,
              'showAttachments': True,
              'mediaInfo': [],
              'relatedRecordsInfo':{
                  'showRelatedRecords':True,
                  'orderByFields':None
              },
              'fieldInfos': getFieldInfos(lyr)
              }
   return popInfo

#get the fields for the popup
def getFieldInfos(layer):
   popInfo=[]
   printMessage("Layer name: " + layer.name)
   #if layer.supports['FEATURECLASS']:
   #     printMessage("Layer has featureclass")
   desc = arcpy.Describe(layer)
   if desc.dataType == "FeatureLayer":
        # Create a fieldinfo object
        allfields= arcpy.ListFields(desc.dataElement.catalogPath)
   else:
        allfields= arcpy.ListFields(desc.catalogPath)
        #return popInfo

   #SmallInteger, Integer, Single, Double, String, Date, OID, Geometry, Blob
   # Iterate through the fields and set them to fieldinfo
   for field in allfields:
        fieldInfos = None
        #printMessage("Field: " + field.name + ":  " + field.type)
        if field.type=='Geometry':
           continue
        if field.type == 'OID':
            oidFldName = field.name
            fieldInfos = {
                'fieldName':field.name,
                'label':field.aliasName,
                'isEditable':field.editable,
                'isEditableOnLayer':field.editable,
                'tooltip':'',
                'visible':False,
                'format':None,
                'stringFieldOption':'textbox'
            }

        elif field.type == 'Integer':
            fieldInfos = {
                'fieldName':field.name,
                'label':field.aliasName,
                'isEditable':field.editable,
                'isEditableOnLayer':field.editable,
                'tooltip':'',
                'visible':True,
                'format':{
                    'places':0,
                    'digitSeparator':True
                },
                'stringFieldOption':'textbox'
            }
        elif field.type == 'Double':
            fieldInfos = {
                'fieldName':field.name,
                'label':field.aliasName,
                'isEditable':field.editable,
                'isEditableOnLayer':field.editable,
                'tooltip':'',
                'visible':True,
                'format':{
                    'places':2,
                    'digitSeparator':True
                    },
                'stringFieldOption':'textbox'
            }
        elif field.type == 'String':
            fieldInfos = {
                'fieldName':field.name,
                'label':field.aliasName,
                'isEditable':field.editable,
                'isEditableOnLayer':field.editable,
                'tooltip':'',
                'visible':True,
                'format':None,
                'stringFieldOption':'textbox'
            }
        else:
            fieldInfos = {
                'fieldName':field.name,
                'label':field.aliasName,
                'isEditable':field.editable,
                'isEditableOnLayer':field.editable,
                'tooltip':'',
                'visible':True,
                'format':None,
                'stringFieldOption':'textbox'
            }
        if fieldInfos is not None:
            popInfo.append(fieldInfos)

   return popInfo

def getFields(layer):
   fields=[]
   desc = arcpy.Describe(layer)
   if hasattr(desc, "layer"):
      catPath = desc.dataElement.catalogPath
   else:
      catPath = desc.catalogPath
   printMessage("Catalog path: "+catPath)
   printMessage(desc.dataType)
   if desc.dataType == "FeatureLayer":
        allfields= arcpy.ListFields(catPath)
   elif desc.dataType == "TableView":
        allfields= arcpy.ListFields(catPath)
   else:
        return fields

   #SmallInteger, Integer, Single, Double, String, Date, OID, Geometry, Blob,Guid
   # Iterate through the fields and set them to fieldinfo
   for field in allfields:
        fieldInfos = None
        #printMessage("Field: " + field.name + ":  " + field.type)
        if field.type=='Geometry':
           continue
        fieldInfos = {
                'alias':field.aliasName,
                'defaultValue':None,
                'domain':None,
                'editable':field.editable,
                'name':field.name,
                'nullable':field.isNullable,
                'sqlType':'sqlTypeOther'
        }
        if field.length:
             fieldInfos['length']=field.length
        #"domain":{"type":"codedValue","name":"cover_type","codedValues":[{"name":"No trees","code":0},{"name":"Trees","code":1}]}

        if field.domain:  #domain contains the name of the domain.  You must look it up to get the full definition using ListDomains
           fieldInfos['domain']={"name":field.domain}
           domains = arcpy.da.ListDomains(getFeatureClassParentWorkspace(catPath))
           for domain in domains:
               if domain.name == field.domain:
                    #printMessage('Domain name: {0}'.format(domain.name))
                    #printMessage('Domain name: {0}'.format(domain.name) )
                    if domain.domainType == 'CodedValue':
                        fieldInfos['domain']['type']='codedValue'
                        codedValuesArray=[]
                        for val,desc in domain.codedValues.iteritems():
                            codedValuesArray.append({"code":val,"name":desc})
                            #[domain.codedValues]
                        fieldInfos['domain']['codedValues'] = codedValuesArray

                        #for val, desc in coded_values.iteritems():
                        #    printMessage('{0} : {1}'.format(val, desc))
                    elif domain.domainType == 'Range':
                        fieldInfos['domain']['type']='range'
                        fieldInfos['domain']['rangeValues']=domain.range
                        #printMessage('Min: {0}'.format(domain.range[0]))
                        #printMessage('Max: {0}'.format(domain.range[1]))

           #for iDomain in arcpy.da.ListDomains(getFeatureClassParentWorkspace(inFeatureClass)):
           #     if iDomain.name == searchDomainName:
           #         return iDomain

        if field.type == 'OID':
            oidFldName = field.name
            fieldInfos['type']='esriFieldTypeOID'
            del fieldInfos['length']
        #elif field.name == 'OBJECTID':
        #    oidFldName = field.name
        #    fieldInfos['type']='esriFieldTypeOID'
        elif field.type == 'Integer':
           fieldInfos['type']='esriFieldTypeInteger'
        elif field.type == 'Single':
           fieldInfos['type']='esriFieldTypeSingle'
           #del fieldInfos['length']
        elif field.type == 'SmallInteger':
           fieldInfos['type']='esriFieldTypeSmallInteger'
           #del fieldInfos['length']
        elif field.type == 'Double':
           fieldInfos['type']='esriFieldTypeDouble'
        elif field.type == 'String':
           fieldInfos['type']='esriFieldTypeString'
        elif field.type == 'Date':
           fieldInfos['type']='esriFieldTypeDate'
        elif field.type == 'Guid':
           fieldInfos['type']='esriFieldTypeGUID'
        elif field.type == 'GlobalID':
           fieldInfos['type']='esriFieldTypeGlobalID'
           fieldInfos['defaultValue']='NEWID() WITH VALUES'
        else:
           printMessage("Unknown field type for " + field.name + ": " + field.type)
           fieldInfos['type']='esriFieldTypeOID'
        if fieldInfos is not None:
            fields.append(fieldInfos)

        #{
        #    "domain": null,
        #    "name": "OBJECTID",
        #    "nullable": false,
        #    "defaultValue": null,
        #    "editable": false,
        #    "alias": "OBJECTID",
        #    "sqlType": "sqlTypeOther",
        #    "type": "esriFieldTypeInteger"
        #},
        #    "domain": null,
        #    "name": "GlobalID",
        #    "nullable": false,
        #    "defaultValue":"NEWID() WITH VALUES",
        #    "editable": false,
        #    "alias": "GlobalID",
        #    "length": 38,
        #    "sqlType": "sqlTypeOther",
        #    "type": "esriFieldTypeGlobalID"

        #if field.name=='OBJECTID':
        #   fieldInfos['type']='esriFieldTypeInteger'

   return fields

#{
#    "id" : <relationshipId1>,
#    "name" : "<relationshipName1>",
#    "relatedTableId" : <relatedTableId1>,
#    "cardinality" : "<esriRelCardinalityOneToOne>|<esriRelCardinalityOneToMany>|<esriRelCardinalityManyToMany>";,//Added at 10.1
#    "role" : "<esriRelRoleOrigin>|<esriRelRoleDestination>";,//Added at 10.1
#    "keyField" : "<keyFieldName2>",//Added at 10.1
#    "composite" : <true>|<false>,//Added at 10.1
#    "relationshipTableId": <attributedRelationshipClassTableId>,  //Added in 10.1. Returned only for attributed relationships
#    "keyFieldInRelationshipTable": "<key field in AttributedRelationshipClass table that matches keyField>" //Added in 10.1. Returned only for attributed relationships
#},

#def getRelationships(lyr,lyrid,cnt,tables,relationshipObj):
def getRelationshipsUnused(lyr,relationshipObj):
   relArr=[]
   desc = arcpy.Describe(lyr)
   if not desc.relationshipClassNames:
      return relArr
   for j,rel in enumerate(desc.relationshipClassNames):
     printMessage("Relationship class name: " + rel)
     relArr.append(relationshipObj[rel])
   return relArr

def getRelationshipsNoGood(lyr,lyrid,cnt,tables,relationships):
   relArr=[]
   desc = arcpy.Describe(lyr)
   if not desc.relationshipClassNames:
      return rel
   if hasattr(desc, "layer"):
         featureName=os.path.basename(desc.layer.catalogPath)
         rootFGDB=desc.layer.catalogPath.replace("\\","/")
   else:
         featureName=os.path.basename(desc.catalogPath)
         rootFGDB=os.path.dirname(desc.catalogPath).replace("\\","/")

   relid=0
   #for index in xrange(0, field_info.count):
   #[u'farm_tracts_inspections__ATTACHREL', u'farm_tractsInspectionRelClass', u'homesites_inspections__ATTACHREL', u'homesitesInspectionRelClass', u'grazing_inspections__ATTACHREL', u'grazing_permitteeRelClass', u'grazing_permitteesInspectionRelClass']
   for j,rel in enumerate(desc.relationshipClassNames):
     relDesc = arcpy.Describe(rootFGDB+"/"+rel)
     if relDesc.isAttachmentRelationship:
          continue
     #printMessage(relDesc)
     #for i in desc:
     #   printMessage(i)

     label = relDesc.backwardPathLabel
     cardinality = relDesc.cardinality
     #key=relDesc.classKey

     originClassKeys=relDesc.originClassKeys
     key=relDesc.destinationClassNames[0]

     printMessage("Origin Class Names")
     printMessage(relDesc.originClassNames)

     printMessage("Origin Class Keys")
     printMessage(relDesc.originClassKeys)

     printMessage("Destination Class Names")
     printMessage(relDesc.destinationClassNames)

     printMessage("Destination Class Keys")
     printMessage(relDesc.destinationClassKeys)

     printMessage("Key type:  "+relDesc.keyType)
     printMessage(relDesc.notification)
     printMessage("backwardPathLabel:  "+relDesc.backwardPathLabel)
     printMessage("forwardPathLabel:  "+relDesc.forwardPathLabel)
     role="esriRelRoleOrigin"
     role="esriRelRoleDestination"

     id=0
     relatedTableId=0
     for i in relDesc.destinationClassNames:
        for j in tables:
          #printMessage(i+":  " + j.datasetName)
          if i == j.datasetName:
            relatedTableId=id
            printMessage("Relationship name: "+i+": " + j.name + "(" + j.datasetName + ") with id: " +str(relatedTableId+cnt))
        id=id+1
     relatedTableId=relatedTableId+cnt

     relObj = {"id":relid,"name":label,"relatedTableId":relatedTableId,"cardinality":"esriRelCardinality"+cardinality,"role":role,"keyField":key,"composite":relDesc.isComposite}
     relArr.append(relObj)
     #relClasses.add(rel)
   return relArr

def getDataIndex(arr,name):
  for j,rel in enumerate(arr):
      #printMessage(str(j) + ": " + rel.name)
      if hasattr(rel, "datasetName"):
         #printMessage(i+":  " + str(j) + ": " + rel.datasetName)
         datasetName=rel.datasetName
      else:
         desc = arcpy.Describe(rel)

         if hasattr(desc, "layer"):
             datasetName=os.path.basename(desc.layer.catalogPath)
         else:
             datasetName=os.path.basename(desc.catalogPath)

      if name == datasetName:
        printMessage("Found relationship name: "+name+": " + rel.name + "(" + datasetName + ") with id: " +str(j))
        return j

  return -1


#"relationships":[{"id":0,"name":"Attributes from homesites","relatedTableId":6,"cardinality":"esriRelCardinalityOneToMany","role":"esriRelRoleOrigin","keyField":"GlobalID","composite":false}],
#"Backward Path Label:", desc.backwardPathLabel
#"Cardinality:", desc.cardinality
#"Class key:", desc.classKey
#"Destination Class Names:", desc.destinationClassNames
#"Forward Path Label:", desc.forwardPathLabel
#"Is Attributed:", desc.isAttributed
#"Is Composite:", desc.isComposite
#"Is Reflexive:", desc.isReflexive
#"Key Type:", desc.keyType
#"Notification Direction:", desc.notification
#"Origin Class Names:", desc.originClassNames
#originClassKeys


# getFeatureClassParentWorkspace: This script gets the geodatabase for a
# feature class. The trick here is that feature classes can be within a
# feature dataset so you need to account for two possible levels in the
# directory structure.
def getFeatureClassParentWorkspace(inFeatureClass):
    describeFC = arcpy.Describe(inFeatureClass)
    if (describeFC.dataType == 'FeatureClass') or (describeFC.dataType == 'Table'):
        workspace1 = describeFC.path
        describeWorkspace1 = arcpy.Describe(workspace1)
        if (describeWorkspace1.dataType == 'FeatureDataset'):
            return describeWorkspace1.path
        return workspace1

    return None

#{"name":"PK__HPL_Coll__F4B70D85F8506C58","fields":"OBJECTID","isAscending":true,"isUnique":true,"description":"clustered, unique, primary key"}


def createReplica(mxd,dataFrame,allData,replicaDestinationPath,toolkitPath,username,serviceName,serverName,minx,miny,maxx,maxy):
  creationDate = time.strftime("%Y-%m-%dT%H:%M:%S")
  sql1=('INSERT INTO GDB_Items("ObjectID", "UUID", "Type", "Name", "PhysicalName", "Path", "Url", "Properties", "Defaults", "DatasetSubtype1", "DatasetSubtype2", "DatasetInfo1", "DatasetInfo2", "Definition", "Documentation", "ItemInfo", "Shape")'
    " select MAX(ObjectID)+1, '{7B6EB064-7BF6-42C8-A116-2E89CD24A000}', '{5B966567-FB87-4DDE-938B-B4B37423539D}', 'MyReplica', 'MYREPLICA', 'MyReplica', '', 1, NULL, NULL, NULL, "
    "'http://"+serverName+"/arcgis/rest/services/"+serviceName+"/FeatureServer', '"+username+"',"
    "'<GPSyncReplica xsi:type=''typens:GPSyncReplica'' xmlns:xsi=''http://www.w3.org/2001/XMLSchema-instance'' xmlns:xs=''http://www.w3.org/2001/XMLSchema'' xmlns:typens=''http://www.esri.com/schemas/ArcGIS/10.3''>"
    "<ReplicaName>MyReplica</ReplicaName><ID>1</ID><ReplicaID>{7b6eb064-7bf6-42c8-a116-2e89cd24a000}</ReplicaID>"
    "<ServiceName>http://"+serverName+"/arcgis/rest/services/"+serviceName+"/FeatureServer</ServiceName>"
    "<Owner>"+username+"</Owner>"
    "<Role>esriReplicaRoleChild</Role><SyncModel>esriSyncModelPerLayer</SyncModel><Direction>esriSyncDirectionBidirectional</Direction><CreationDate>"+creationDate+"</CreationDate><LastSyncDate>1970-01-01T00:00:01</LastSyncDate>"
    "<ReturnsAttachments>true</ReturnsAttachments><SpatialRelation>esriSpatialRelIntersects</SpatialRelation><QueryGeometry xsi:type=''typens:PolygonN''><HasID>false</HasID><HasZ>false</HasZ><HasM>false</HasM><Extent xsi:type=''typens:EnvelopeN''>"
    "<XMin>"+minx+"</XMin><YMin>"+miny+"</YMin><XMax>"+maxx+"</XMax><YMax>"+maxy+"</YMax></Extent><RingArray xsi:type=''typens:ArrayOfRing''><Ring xsi:type=''typens:Ring''>"
    "<PointArray xsi:type=''typens:ArrayOfPoint''>"
    "<Point xsi:type=''typens:PointN''><X>"+minx+"</X><Y>"+miny+"</Y></Point><Point xsi:type=''typens:PointN''><X>"+maxx+"</X><Y>"+miny+"</Y></Point>"
    "<Point xsi:type=''typens:PointN''><X>"+maxx+"</X><Y>"+maxy+"</Y></Point><Point xsi:type=''typens:PointN''><X>"+minx+"</X><Y>"+maxy+"</Y></Point>"
    "<Point xsi:type=''typens:PointN''><X>"+minx+"</X><Y>"+miny+"</Y></Point></PointArray></Ring></RingArray>"
    "<SpatialReference xsi:type=''typens:ProjectedCoordinateSystem''><WKT>PROJCS[&quot;WGS_1984_Web_Mercator_Auxiliary_Sphere&quot;,GEOGCS[&quot;GCS_WGS_1984&quot;,DATUM[&quot;D_WGS_1984&quot;,SPHEROID[&quot;WGS_1984&quot;,6378137.0,298.257223563]],PRIMEM[&quot;Greenwich&quot;,0.0],UNIT[&quot;Degree&quot;,0.0174532925199433]],PROJECTION[&quot;Mercator_Auxiliary_Sphere&quot;],PARAMETER[&quot;False_Easting&quot;,0.0],PARAMETER[&quot;False_Northing&quot;,0.0],PARAMETER[&quot;Central_Meridian&quot;,0.0],PARAMETER[&quot;Standard_Parallel_1&quot;,0.0],PARAMETER[&quot;Auxiliary_Sphere_Type&quot;,0.0],UNIT[&quot;Meter&quot;,1.0],AUTHORITY[&quot;EPSG&quot;,3857]]</WKT>"
    "<XOrigin>-20037700</XOrigin><YOrigin>-30241100</YOrigin><XYScale>10000</XYScale><ZOrigin>-100000</ZOrigin><ZScale>10000</ZScale><MOrigin>-100000</MOrigin><MScale>10000</MScale><XYTolerance>0.001</XYTolerance><ZTolerance>0.001</ZTolerance><MTolerance>0.001</MTolerance><HighPrecision>true</HighPrecision><WKID>102100</WKID><LatestWKID>3857</LatestWKID></SpatialReference></QueryGeometry>")

  #"<GPSyncDatasets xsi:type=''typens:ArrayOfGPSyncDataset''>"
  #"<GPSyncDataset xsi:type=''typens:GPSyncDataset''><DatasetID>5</DatasetID><DatasetName>"+layerName+"</DatasetName><DatasetType>esriDTFeatureClass</DatasetType>"
  #"<LayerID>"+str(id)+"</LayerID><LayerName>"+layerName+"</LayerName><Direction>esriSyncDirectionBidirectional</Direction><ReplicaServerGen xsi:type=''xs:long''>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type=''xs:long''>1000</ReplicaClientDownloadGen>"
  #"<ReplicaClientUploadGen xsi:type=''xs:long''>1000</ReplicaClientUploadGen><ReplicaClientAcknowledgeUploadGen xsi:type=''xs:long''>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>true</IncludeRelated>"
  #"<QueryOption>esriRowsTypeFilter</QueryOption><HasAttachments>true</HasAttachments></GPSyncDataset>"
  #
  #"<GPSyncDataset xsi:type=''typens:GPSyncDataset''><DatasetID>6</DatasetID><DatasetName>"+layerName+"__ATTACH</DatasetName><DatasetType>esriDTTable</DatasetType>"
  #"<LayerID>0</LayerID>"
  #"<LayerName>"+layerName+"</LayerName><Direction>esriSyncDirectionBidirectional</Direction>"
  #"<ReplicaServerGen xsi:type=''xs:long''>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type=''xs:long''>1000</ReplicaClientDownloadGen><ReplicaClientUploadGen xsi:type=''xs:long''>1000</ReplicaClientUploadGen>"
  #"<ReplicaClientAcknowledgeUploadGen xsi:type=''xs:long''>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>false</IncludeRelated><QueryOption>esriRowsTypeFilter</QueryOption><IsAttachment>true</IsAttachment></GPSyncDataset>"

  #lyrs=[]
  #for lyr in arcpy.mapping.ListLayers(mxd, "", dataFrame):
  #   # Exit if the current layer is not a service layer.
  #   if lyr.isServiceLayer or lyr.supports("SERVICEPROPERTIES"):  # or not lyr.visible
  #      continue
  #   lyrs.append(lyr)
  #for lyr in arcpy.mapping.ListTableViews(mxd, "", dataFrame):
  #   lyrs.append(lyr)

  id=0
  sql2=[]
  sql3=[]
  for lyr in allData:
     desc = arcpy.Describe(lyr)
     if hasattr(desc, "layer"):
         featureName=os.path.basename(desc.layer.catalogPath)
         inFeaturesGDB=desc.layer.catalogPath.replace("\\","/")
     else:
         featureName=os.path.basename(desc.catalogPath)
         inFeaturesGDB=os.path.dirname(desc.catalogPath).replace("\\","/")

     hasAttachments="false"
     if arcpy.Exists(inFeaturesGDB+"/"+featureName+"__ATTACH"):
         hasAttachments="true"
     sql1=sql1+ ("<GPSyncDatasets xsi:type=''typens:ArrayOfGPSyncDataset''>"
        "<GPSyncDataset xsi:type=''typens:GPSyncDataset''><DatasetID>"+str(id)+"</DatasetID><DatasetName>"+lyr.name+"</DatasetName><DatasetType>esriDTFeatureClass</DatasetType>"
        "<LayerID>"+str(id)+"</LayerID><LayerName>"+lyr.name+"</LayerName><Direction>esriSyncDirectionBidirectional</Direction><ReplicaServerGen xsi:type=''xs:long''>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type=''xs:long''>1000</ReplicaClientDownloadGen>"
        "<ReplicaClientUploadGen xsi:type=''xs:long''>1000</ReplicaClientUploadGen><ReplicaClientAcknowledgeUploadGen xsi:type=''xs:long''>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>true</IncludeRelated>"
        "<QueryOption>esriRowsTypeFilter</QueryOption><HasAttachments>true</HasAttachments></GPSyncDataset>")

     sql2.append(('INSERT INTO GDB_Items("ObjectID", "UUID", "Type", "Name", "PhysicalName", "Path", "Url", "Properties", "Defaults", "DatasetSubtype1", "DatasetSubtype2", "DatasetInfo1", "DatasetInfo2", "Definition", "Documentation", "ItemInfo", "Shape")'
        " select MAX(ObjectID)+1, '{AE8D3C7E-9890-4BF4-B946-5BE50A1CC2"+str(format(id, '02'))+"}', '{D86502F9-9758-45C6-9D23-6DD1A0107B47}', '"+featureName+"', '"+featureName.upper()+"', 'MyReplica\\"+featureName+"', '', 1, NULL, NULL, NULL, NULL, NULL, "
        "'<GPSyncDataset xsi:type=''typens:GPSyncDataset'' xmlns:xsi=''http://www.w3.org/2001/XMLSchema-instance'' xmlns:xs=''http://www.w3.org/2001/XMLSchema'' xmlns:typens=''http://www.esri.com/schemas/ArcGIS/10.3''>"
        "<DatasetID>"+str(id)+"</DatasetID>"
        "<DatasetName>"+lyr.name+"</DatasetName>"
        "<DatasetType>esriDTFeatureClass</DatasetType><LayerID>"+str(id)+"</LayerID><LayerName>"+lyr.name+"</LayerName><Direction>esriSyncDirectionBidirectional</Direction>"
        "<ReplicaServerGen xsi:type=''xs:long''>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type=''xs:long''>1000</ReplicaClientDownloadGen><ReplicaClientUploadGen xsi:type=''xs:long''>1000</ReplicaClientUploadGen>"
        "<ReplicaClientAcknowledgeUploadGen xsi:type=''xs:long''>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>true</IncludeRelated><QueryOption>esriRowsTypeFilter</QueryOption>"
        "<HasAttachments>"+hasAttachments+"</HasAttachments></GPSyncDataset>', NULL, NULL, NULL from GDB_Items"))

     printMessage("Loading " + lyr.name)
     if arcpy.Exists(inFeaturesGDB+"/"+featureName+"__ATTACH"):
        printMessage("Found attachment table: " + featureName+"__ATTACH")
        sql1=sql1+ ("<GPSyncDatasets xsi:type=''typens:ArrayOfGPSyncDataset''>"
        "<GPSyncDataset xsi:type=''typens:GPSyncDataset''><DatasetID>"+str(id)+"</DatasetID><DatasetName>"+lyr.name+"</DatasetName><DatasetType>esriDTFeatureClass</DatasetType>"
        "<LayerID>"+str(id)+"</LayerID><LayerName>"+lyr.name+"</LayerName><Direction>esriSyncDirectionBidirectional</Direction><ReplicaServerGen xsi:type=''xs:long''>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type=''xs:long''>1000</ReplicaClientDownloadGen>"
        "<ReplicaClientUploadGen xsi:type=''xs:long''>1000</ReplicaClientUploadGen><ReplicaClientAcknowledgeUploadGen xsi:type=''xs:long''>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>true</IncludeRelated>"
        "<QueryOption>esriRowsTypeFilter</QueryOption><HasAttachments>true</HasAttachments></GPSyncDataset>")

        sql3.append(('INSERT INTO GDB_Items("ObjectID", "UUID", "Type", "Name", "PhysicalName", "Path", "Url", "Properties", "Defaults", "DatasetSubtype1", "DatasetSubtype2", "DatasetInfo1", "DatasetInfo2", "Definition", "Documentation", "ItemInfo", "Shape")'
        " select MAX(ObjectID)+1, '{55C5E7E4-834D-4D44-A12C-991E7F8B46"+str(format(id, '02'))+"}', '{D86502F9-9758-45C6-9D23-6DD1A0107B47}', '"+featureName+"__ATTACH', '"+featureName.upper()+"__ATTACH', 'MyReplica\\"+featureName+"__ATTACH', '', 1, NULL, NULL, NULL, NULL, NULL, "
        "'<GPSyncDataset xsi:type=''typens:GPSyncDataset'' xmlns:xsi=''http://www.w3.org/2001/XMLSchema-instance'' xmlns:xs=''http://www.w3.org/2001/XMLSchema'' xmlns:typens=''http://www.esri.com/schemas/ArcGIS/10.3''>"
        "<DatasetID>"+str(id)+"</DatasetID>"
        "<DatasetName>"+featureName+"__ATTACH</DatasetName><DatasetType>esriDTTable</DatasetType><LayerID>"+str(id)+"</LayerID><LayerName>"+featureName+"</LayerName><Direction>esriSyncDirectionBidirectional</Direction>"
        "<ReplicaServerGen xsi:type=''xs:long''>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type=''xs:long''>1000</ReplicaClientDownloadGen><ReplicaClientUploadGen xsi:type=''xs:long''>1000</ReplicaClientUploadGen>"
        "<ReplicaClientAcknowledgeUploadGen xsi:type=''xs:long''>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>false</IncludeRelated><QueryOption>esriRowsTypeFilter</QueryOption>"
        "<IsAttachment>true</IsAttachment></GPSyncDataset>',"
        " NULL, NULL, NULL from GDB_Items"))

     id = id + 1

  #sql3=('INSERT INTO GDB_Items("ObjectID", "UUID", "Type", "Name", "PhysicalName", "Path", "Url", "Properties", "Defaults", "DatasetSubtype1", "DatasetSubtype2", "DatasetInfo1", "DatasetInfo2", "Definition", "Documentation", "ItemInfo", "Shape")'
  #" select MAX(ObjectID)+1, '{55C5E7E4-834D-4D44-A12C-991E7F8B4645}', '{D86502F9-9758-45C6-9D23-6DD1A0107B47}', '"+layerName+"__ATTACH', '"+layerName.upper()+"__ATTACH', 'MyReplica_"+str(id)+"\\"+layerName+"__ATTACH', '', 1, NULL, NULL, NULL, NULL, NULL, "
  #"'<GPSyncDataset xsi:type=''typens:GPSyncDataset'' xmlns:xsi=''http://www.w3.org/2001/XMLSchema-instance'' xmlns:xs=''http://www.w3.org/2001/XMLSchema'' xmlns:typens=''http://www.esri.com/schemas/ArcGIS/10.3''>"
  #"<DatasetID>6</DatasetID>"
  #"<DatasetName>"+layerName+"__ATTACH</DatasetName><DatasetType>esriDTTable</DatasetType><LayerID>0</LayerID><LayerName>"+layerName+"</LayerName><Direction>esriSyncDirectionBidirectional</Direction>"
  #"<ReplicaServerGen xsi:type=''xs:long''>2590</ReplicaServerGen><ReplicaClientDownloadGen xsi:type=''xs:long''>1000</ReplicaClientDownloadGen><ReplicaClientUploadGen xsi:type=''xs:long''>1000</ReplicaClientUploadGen>"
  #"<ReplicaClientAcknowledgeUploadGen xsi:type=''xs:long''>1000</ReplicaClientAcknowledgeUploadGen><UseGeometry>true</UseGeometry><IncludeRelated>false</IncludeRelated><QueryOption>esriRowsTypeFilter</QueryOption>"
  #"<IsAttachment>true</IsAttachment></GPSyncDataset>',"
  #" NULL, NULL, NULL from GDB_Items")

  sql4='update "GDB_ServiceItems" set "ItemInfo"=replace("ItemInfo",|"capabilities":"Query"|,|"capabilities":"Create,Delete,Query,Update,Editing,Sync"|)'
  sql4=sql4.replace("|","'")
  #sql5='update "GDB_Items" set ObjectId=ROWID'
  sql1=sql1+("</GPSyncDatasets></GPSyncReplica>'"
   ", NULL, NULL, NULL from GDB_Items")
  #serviceRep=[sql1,sql2,sql4]

  arcpy.CreateRuntimeContent_management(mxd.filePath,
              replicaDestinationPath + os.sep + serviceName,
              serviceName,"#","#",
              "FEATURE_AND_TABULAR_DATA","NON_OPTIMIZE_SIZE","ONLINE","PNG","1","#")
  filenames = next(os.walk(replicaDestinationPath + "/"+serviceName+"/data/"))[2]
  printMessage("Rename " + replicaDestinationPath + "/"+serviceName+"/data/"+filenames[0]+" to "+ replicaDestinationPath+"/"+serviceName+".geodatabase")
  #if offline geodatabase exists, must delete first
  newFullReplicaDB=replicaDestinationPath+"/"+serviceName+".geodatabase"
  try:
     if os.path.exists(newFullReplicaDB):
        os.rmdir(newFullReplicaDB)
  except:
     printMessage("Unable to remove old replica geodatabase")

  os.rename(replicaDestinationPath + "/"+serviceName+"/data/"+filenames[0], newFullReplicaDB)
  try:
     os.rmdir(replicaDestinationPath + "/"+serviceName+"/data/")
     os.rmdir(replicaDestinationPath + "/"+serviceName)
  except:
     printMessage("Unable to remove replica folders")

  name=replicaDestinationPath + "/"+serviceName+".sql"
  with open(name,'w') as f:
       f.write("SELECT load_extension( 'stgeometry_sqlite.dll', 'SDE_SQL_funcs_init');\n")
       f.write(sql1)
       f.write(";\n")

       for i in sql2:
          f.write(i)
          f.write(";\n")
       for i in sql3:
          f.write(i)
          f.write(";\n")
       f.write(sql4)
       f.write(";\n")
       f.close()
  #printMessage("Running \"" + toolkitPath+"/spatialite/spatialite.exe\" \"" + newFullReplicaDB + "\"  < " + name)
  printMessage("Running \"" + toolkitPath+"/bin/sqlite3.exe\" \"" + newFullReplicaDB + "\"  < " + name)

  #try:
  #   os.system(toolkitPath+"/bin/sqlite3.exe \"" + newFullReplicaDB + "\"  < \"" + name + "\" >>" + replicaDestinationPath + os.sep + "output.log 2>&1")
  #except:
  #   printMessage("Unable to run sql commands")
  #

def getIndexes(lyr):
   indexes=[]
   desc = arcpy.Describe(lyr)
   if desc.dataType == "FeatureLayer":
       lyrindexes = arcpy.ListIndexes(desc.dataElement.catalogPath)
   elif desc.dataType == "TableView":
       lyrindexes = arcpy.ListIndexes(desc.catalogPath)

   for index in lyrindexes:
       fields=[]
       for field in index.fields:
          fields.append(field.name)

       indexes.append({"name":index.name,"fields":",".join(fields),"isAscending":index.isAscending,"isUnique":index.isUnique})
   return indexes

   #printMessage("Name        : {0}".format(index.name))
   #printMessage("IsAscending : {0}".format(index.isAscending))
   #printMessage("IsUnique    : {0}".format(index.isUnique))

def hasAttachments(tbl):
    return arcpy.Exists(tbl+"__ATTACH")

def hasEditorTracking(tbl):
    desc = arcpy.Describe(tbl)
    return desc.editorTrackingEnabled

def num(s):
    try:
        return int(s)
    except ValueError:
        return float(s)

def getSymbolColor(sym):
    printMessage("here")

def getPointSymbol(sym):
    obj = {}
    obj['type']="esriSMS"
    obj['style']="esriSMSCircle"
    obj['size']= 4
    obj['angle']= 0
    obj['xoffset']= 0
    obj['yoffset']=  0
    obj['outline']={}
    obj['outline']['width']= 1
    symb = sym.getElementsByTagName("CIMSymbolLayer")[0]
    size = symb.getElementsByTagName("Size")
    width = symb.getElementsByTagName("Width")
    if len(size) > 0:
       printMessage("       Size: " + size[0].childNodes[0].nodeValue)
       obj['size']= num(size[0].childNodes[0].nodeValue)
    if len(width)>0:
       printMessage("       Width: " + width[0].childNodes[0].nodeValue)
       obj['width']=num(width[0].childNodes[0].nodeValue)  
    obj = getSymbolColor(sym,obj)               
    return obj
    
def getPolygonSymbol(sym):
    obj = {}
    obj['type']="esriSFS"
    obj['style']="esriSFSSolid"
    obj['outline']={}
    symb = sym.getElementsByTagName("CIMSymbolLayer")[0]
    if symb.getAttribute("xsi:type")=="typens:CIMFill":
                    obj['outline']['type']="esriSLS"
                    #2drawingInfo['renderer']['symbol']['outline']['style']="esriSFSSolid"
                    obj['outline']['style']="esriSLSSolid"
                    
    else:
                    obj['outline']['type']="esriSLS"
                    obj['outline']['style']="esriSLSSolid"
    obj = getSymbolColor(sym,obj)  
    return obj

def getSymbolColor(sym,obj):
    patt = sym.getElementsByTagName("Pattern")
    if len(patt)>0:
         color = patt[0].getElementsByTagName("Color")
         colorStr = str(color[0].getElementsByTagName("R")[0].childNodes[0].nodeValue) + "," + str(color[0].getElementsByTagName("G")[0].childNodes[0].nodeValue) + "," + str(color[0].getElementsByTagName("B")[0].childNodes[0].nodeValue)
         if patt[0].getAttribute("xsi:type")=="typens:CIMFilledStroke":
            obj['outline']['color']=[ int(color[0].getElementsByTagName("R")[0].childNodes[0].nodeValue), int(color[0].getElementsByTagName("G")[0].childNodes[0].nodeValue), int(color[0].getElementsByTagName("B")[0].childNodes[0].nodeValue),255]
         else:
            obj['color']=[ int(color[0].getElementsByTagName("R")[0].childNodes[0].nodeValue), int(color[0].getElementsByTagName("G")[0].childNodes[0].nodeValue), int(color[0].getElementsByTagName("B")[0].childNodes[0].nodeValue),255]
         printMessage("         Color (polygon): " + colorStr)
    return obj

def getSymbolLayers(sym):
    return {}
    #for l in sym.childNodes:

def getSymbolLayersa(sym):
    for l in sym.childNodes:
       printMessage("   " + l.tagName + ": " + l.getAttribute("xsi:type"))
       if l.tagName=='SymbolLayers':
          #printMessage("    " + l.tagName)
          for m in l.childNodes:
              printMessage("     " + m.tagName + ": " + m.getAttribute("xsi:type"))
              if m.tagName == "CIMSymbolLayer":
                 if m.getAttribute("xsi:type")=="typens:CIMFill":
                    drawingInfo['renderer']['symbol']['outline']['type']="esriSLS"
                    #2drawingInfo['renderer']['symbol']['outline']['style']="esriSFSSolid"
                    drawingInfo['renderer']['symbol']['outline']['style']="esriSLSSolid"
                    
                 else:
                    drawingInfo['renderer']['symbol']['outline']['type']="esriSLS"
                    drawingInfo['renderer']['symbol']['outline']['style']="esriSLSSolid"
                     
                    #{"type":"esriSLS",
                    #"style":"esriSLSSolid",

                 for n in m.childNodes:
                    printMessage("      " + n.tagName + ": " + n.getAttribute("xsi:type"))
                    if n.tagName=='Size':
                      printMessage("       Size: " + n.childNodes[0].nodeValue)
                      drawingInfo['renderer']['symbol']['size']= num(n.childNodes[0].nodeValue)
                    if n.tagName=='Width':
                      printMessage("       Width: " + n.childNodes[0].nodeValue)
                      drawingInfo['renderer']['symbol']['width']=num(n.childNodes[0].nodeValue)
                    
                    if n.tagName=='Pattern':
                      color = n.getElementsByTagName("Color")
                      colorStr = str(color[0].getElementsByTagName("R")[0].childNodes[0].nodeValue) + "," + str(color[0].getElementsByTagName("G")[0].childNodes[0].nodeValue) + "," + str(color[0].getElementsByTagName("B")[0].childNodes[0].nodeValue)
                      if m.getAttribute("xsi:type")=="typens:CIMFilledStroke":
                         drawingInfo['renderer']['symbol']['outline']['color']=[ int(color[0].getElementsByTagName("R")[0].childNodes[0].nodeValue), int(color[0].getElementsByTagName("G")[0].childNodes[0].nodeValue), int(color[0].getElementsByTagName("B")[0].childNodes[0].nodeValue),255]
                      else:
                         drawingInfo['renderer']['symbol']['color']=[ int(color[0].getElementsByTagName("R")[0].childNodes[0].nodeValue), int(color[0].getElementsByTagName("G")[0].childNodes[0].nodeValue), int(color[0].getElementsByTagName("B")[0].childNodes[0].nodeValue),255]
                      printMessage("         Color (polygon): " + colorStr)
                    
                    if n.tagName=='Symbol':
                      for o in n.childNodes:
                        printMessage("       " + o.tagName + ": " + o.getAttribute("xsi:type"))
                        if o.tagName=='SymbolLayers':
                          for p in o.childNodes:
                             printMessage("        " + p.tagName + ": " + p.getAttribute("xsi:type"))
                             if p.tagName =='CIMSymbolLayer':
                                for q in p.childNodes:
                                   printMessage("         " + q.tagName + ": " + q.getAttribute("xsi:type"))
                                   if q.tagName=='Pattern':
                                      color = q.getElementsByTagName("Color")
                                      colorStr = str(color[0].getElementsByTagName("R")[0].childNodes[0].nodeValue) + "," + str(color[0].getElementsByTagName("G")[0].childNodes[0].nodeValue) + "," + str(color[0].getElementsByTagName("B")[0].childNodes[0].nodeValue)
                                      drawingInfo['renderer']['symbol']['color']=[ int(color[0].getElementsByTagName("R")[0].childNodes[0].nodeValue), int(color[0].getElementsByTagName("G")[0].childNodes[0].nodeValue), int(color[0].getElementsByTagName("B")[0].childNodes[0].nodeValue), 255]
                                      printMessage("          Color (point): " + colorStr)


#see http://resources.arcgis.com/en/help/arcgis-rest-api/index.html#//02r30000019t000000
# and http://resources.arcgis.com/en/help/arcgis-rest-api/index.html#//02r3000000n5000000
#"symbol":{ "type": "esriSMS", "style": "esriSMSSquare", "color": [76,115,0,255], "size": 8, "angle": 0, "xoffset": 0, "yoffset": 0, "outline": { "color": [152,230,0,255], "width": 1 } }
def getSymbol(lyr,sym,name):
   drawingInfo= {
     "renderer": {
      "label": "",
      "description": ""
     },
     "transparency": lyr.transparency,
     "labelingInfo": None
   }

   if sym[0].getAttribute("xsi:type") == "typens:CIMUniqueValueSymbolizer":
         drawingInfo['renderer']['type']="uniqueValue"
         drawingInfo['renderer']['uniqueValueInfos']=[]
   else:
         drawingInfo['renderer']['type']="simple"
         #drawingInfo['renderer']['symbol']={}
         #drawingInfo['renderer']['symbol']['outline']={}

   #renderer->uniqueValueInfos
   printMessage("******Creating symbology for " + name + "*******")
   for i in sym:
      printMessage(i.tagName + ": " + i.getAttribute("xsi:type"))
      #printMessage(i)
      #printMessage(str(i.childNodes.length))

      for j in i.childNodes:
         printMessage(" " + j.tagName + ": " + j.getAttribute("xsi:type"))
         if i.tagName=='Fields':
            c=1
            for k in j.childNodes:
               drawingInfo['renderer']['field'+str(c)]=j.childNodes[k].nodeValue
               c=c+1

         #get first symbol
         if j.tagName=='Symbol':
            for k in j.childNodes:
               printMessage("  " + k.tagName + ": " + k.getAttribute("xsi:type"))
               #get second symbol
               #if k.getAttribute("xsi:type")=="typens:CIMSymbolReference":
               if k.tagName=='Symbol':
                 if k.getAttribute("xsi:type")=="typens:CIMPointSymbol":
                    drawingInfo['renderer']['symbol'] = getPointSymbol(k)
                 elif k.getAttribute("xsi:type")=="typens:CIMPolygonSymbol":
                     drawingInfo['renderer']['symbol']=getPolygonSymbol(k)
               elif k.tagName=='SymbolLayers':
                    #drawingInfo['renderer']['symbol'] = getSymbolLayers(k)
                    drawingInfo['renderer']['uniqueValueInfos']=getSymbolLayers(k)
   
   return drawingInfo

   #if sym[0].getElementsByTagName("OutlineColor").length > 0:
   #  printMessage("OutlineColor len: "+ str(sym[0].getElementsByTagName("OutlineColor").length))

   #printMessage("Lookup Symbol length: " + name + ": length: " + str(sym.length))
   #printMessage("Type: "+sym[0].getAttribute("xsi:type"))
   #printMessage("Symbol children length: " + str(sym[0].childNodes.length))

   #symbolizer = sym.getElementsByTagName("Symbolizer")
   #if sym.length > 0:
   #    symtags = sym[0].getElementsByTagName("Symbol")
   #    if symtags[0].getAttribute("xsi:type") == 'typens:CIMSymbolReference':
   #        printMessage("Found new symbol")
   #        symtags = symtags[0].getElementsByTagName("Symbol")

       #if symtags.length > 0:
   #    for symbol in symtags:
   #        symbol = symtags[0].getElementsByTagName("Symbol")
   #        if symbol.length > 0:
   #           printMessage("Symbol len: "+ str(symbol.length))
   #        else:
   #           continue

   #        printMessage("Symbol: "+symbol[0].getAttribute("xsi:type"))
   #        symbollayers = symbol[0].getElementsByTagName("SymbolLayers")
   #        printMessage("SymbolLayers: "+symbollayers[0].getAttribute("xsi:type"))

   #        symbollayer  = symbollayers[0].getElementsByTagName("CIMSymbolLayer")
   #        printMessage("CIMSymbolLayer 1: "+str(symbollayer.length) + " - " + symbollayer[0].getAttribute("xsi:type"))

   #        fillcolor = symbollayer[0].getElementsByTagName("FillColor")
   #        if fillcolor.length>0:
   #            printMessage("FillColor->R,G,B: " + fillcolor[0].childNodes[0].nodeValue)

    #       outlinecolor = symbollayer[0].getElementsByTagName("OutlineColor")
    #       if outlinecolor.length>0:
    #           printMessage("OutlineColor->R,G,B: " + outlinecolor[0].childNodes[0].nodeValue)

     #      size = symbollayer[0].getElementsByTagName("Size") #[0].childNodes[0].nodeValue
     #      if size.length>0:
     #          printMessage("Size: " + size[0].childNodes[0].nodeValue)

     #      type = symbollayer[0].getElementsByTagName("Type")  #[0].childNodes[0].nodeValue
     #      if type.length>0:
     #          printMessage("Type: " + type[0].childNodes[0].nodeValue)

           #if symbollayer[0].getElementsByTagName("FillColor").length > 0:
           #   fill = symbollayer[0].getElementsByTagName("FillColor")[0].childNodes[0].nodeValue
           #   outline = symbollayer[0].getElementsByTagName("FillColor")[0].childNodes[0].nodeValue
           #   size = symbollayer[0].getElementsByTagName("FillColor")[0].childNodes[0].nodeValue
           #   printMessage("FillColor->R,G,B: " + str(fill))
           #   printMessage("OutlineColor->R,G,B: " + str(outline))
           #   printMessage("Size: "+str(size))
           #elif symbollayer[0].getElementsByTagName("FillColor").length > 0:
           #   fill = symbollayer[0].getElementsByTagName("FillColor")[0].childNodes[0].nodeValue
           #   outline = symbollayer[0].getElementsByTagName("FillColor")[0].childNodes[0].nodeValue
           #   size = symbollayer[0].getElementsByTagName("FillColor")[0].childNodes[0].nodeValue
           #   printMessage("FillColor->R,G,B: " + str(fill))
           #   printMessage("OutlineColor->R,G,B: " + str(outline))
           #   printMessage("Size: "+str(size))
           #else:
           #   color = symbollayer[0].getElementsByTagName("Color")
           #   printMessage("Color: "+color[0].getAttribute("xsi:type"))
           #   #[0].childNodes[0].nodeValue
           #   colorStr = str(color[0].getElementsByTagName("R")[0].childNodes[0].nodeValue) + "," + str(color[0].getElementsByTagName("G")[0].childNodes[0].nodeValue) + "," + str(color[0].getElementsByTagName("B")[0].childNodes[0].nodeValue)
           #   printMessage("Color values: " + colorStr)


   #        printMessage("")


   #return drawingInfo

  #printMessage(lyr.name + ": " + lyr.symbologyType)
  #if lyr.symbologyType!="OTHER":
  #   for i in lyr.symbology:
  #      printMessage(i)

#Points
#"drawingInfo": {
#  "renderer": {
#   "type": "simple",
#   "symbol": {
#    "type": "esriSMS",
#    "style": "esriSMSCircle",
#    "color": [135,99,0,255],
#    "size": 4,
#    "angle": 0,
#    "xoffset": 0,
#    "yoffset": 0,
#    "outline": {
#    "color": [0,0,0,255],
#    "width": 1
#    }
#   },
#   "label": "",
#   "description": ""
#  },
#  "transparency": 0,
#  "labelingInfo": null
# },

#Polygons
#"drawingInfo":
#{
#"renderer":
#{
#"type":"simple",
#"symbol":
#{
#"type":"esriSFS",
#"style":"esriSFSSolid",
#"color":[76,129,205,191],
#"outline":
#{
#"type":"esriSLS",
#"style":"esriSLSSolid",
#"color":[0,0,0,255],
#"width":0.75
#}
#}
#},
#"transparency":0,"labelingInfo":null
#},


#symbol={"angle":0}
#if lyr.symbologyType == "GRADUATED_COLORS":
#     lyr.symbology.valueField = "POP2000"
#     lyr.symbology.numClasses = 4
#     lyr.symbology.classBreakValues = [250000, 999999, 4999999, 9999999, 35000000]
#     lyr.symbology.classBreakLabels = ["250,000 to 999,999", "1,000,000 to 4,999,999",
#                                       "5,000,000 to 9,999,999", "10,000,000 to 35,000,000"]

#return ""































class LayerExtras(object):
    """ An object to hold attributes loaded from xml inside the msd."""

    name = ""
    symbologyFieldName = ""


class MxdExtras(dict):
    """ Exposes extra MXD details by raiding an exported msd

        Treat this object as a dictionary with layer name as the key and a custom object
        with desired attributes as the value.
        You must have write access to MXD directory (creates temporary msd file).
        Only layers in the first dataframe are accessed.

    """

    LYR_NAME_NODE = "Name"
    LYR_SYMBOL_NODE = "Symbolizer"
    LYR_FIELD_NODE = "Field"
    MSD_SUFFIX = "_MxdExtrasTemp.msd"
    MXD_SUFFIX = ".mxd"
    EXCLUDED_FILE_NAMES = ["DocumentInfo.xml", "layers/layers.xml"]
    mxdPath = ""

    def __init__(self, mxdPath):

        self.loadMxdPath(mxdPath)


    def loadMxdPath(self, mxd):
        """ Load mxd from file path """

        #self.mxdPath = mxdPath.lower()
        #mxd = mapping.MapDocument(self.mxdPath)
        msdPath = os.path.abspath(os.path.dirname(__file__)).replace("\\","/")+"/output.msd"

        #msdPath = self.mxdPath.replace(self.MXD_SUFFIX, self.MSD_SUFFIX)

        # Delete temporary msd if it exists
        if os.path.exists(msdPath):
            os.remove(msdPath)

        mapping.ConvertToMSD(mxd,msdPath)

        zz = zipfile.ZipFile(msdPath)

        for fileName in (fileName for fileName in zz.namelist() if not fileName in self.EXCLUDED_FILE_NAMES):
            printMessage("Opening: " + fileName)
            dom = parse(zz.open(fileName))
            symb = dom.getElementsByTagName(self.LYR_SYMBOL_NODE)
            #printMessage("Found: " + str(symb.length))
            name,lyr = self.loadMsdLayerDom(dom)
            if name != "":
               self[name] = lyr
        del zz
        #os.remove(msdPath)

    def loadMsdLayerDom(self, dom):
        """ Load dom created from xml file inside the msd. """

        lyr = LayerExtras()

        # Layer name
        layer = dom.getElementsByTagName(self.LYR_NAME_NODE)
        if layer.length==0:
           printMessage("Layer not found")
           return "",""
        lyr.name = layer[0].childNodes[0].nodeValue

        labelingInfo=None
        label = dom.getElementsByTagName("LabelClasses")
        if label.length>0:
           cim = label[0].getElementsByTagName("CIMLabelClass")
           if cim.length > 0:
              expr = cim[0].getElementsByTagName("Expression")
              if expr.length > 0:
                 expression = expr[0].childNodes[0].nodeValue
              labelingInfo = {
                "labelPlacement":"esriServerPointLabelPlacementAboveRight",
                "where":null,
                "labelExpression":expression,
                "useCodedValues":true,
                "symbol": {
                   "type":"esriTS",
                   "color":[0,0,0,255],
                   "backgroundColor":null,
                   "borderLineColor":null,
                   "borderLineSize":null,
                   "verticalAlignment":"bottom",
                   "horizontalAlignment":"center",
                   "rightToLeft":false,
                   "angle":0,
                   "xoffset":0,
                   "yoffset":0,
                   "kerning":true,
                   "haloColor":null,
                   "haloSize":null,
                   "font":{"family":"Arial","size":8,"style":"normal","weight":"normal","decoration":"none"}
                },
                "minScale":0,
                "maxScale":0
              }

        # Symbology field name.  Create a symbol object
        drawingInfo= {
         "renderer": {
          "type": "simple",
          "symbol": {
           "type": "esriSMS",
           "style": "esriSMSCircle",
           "color": [135,99,0,255],
           "size": 4,
           "angle": 0,
           "xoffset": 0,
           "yoffset": 0,
           "outline": {
           "color": [0,0,0,255],
           "width": 1
           }
          },
          "label": "",
          "description": ""
         },
         "transparency": 0,
         "labelingInfo": labelingInfo
        }

        symbolizer = dom.getElementsByTagName("Symbolizer")
        if symbolizer.length > 0:
           symbologyElement = symbolizer[0]
           symtags = symbologyElement.getElementsByTagName("Symbol")
           if symtags.length > 0:
               symbol = symtags[0].getElementsByTagName("Symbol")[0]
               #lyr.symbologyFieldName = symtags[0].childNodes[0].nodeValue
               if symbol.getElementsByTagName("FillColor").length > 0:
                  fill = symbol.getElementsByTagName("FillColor")[0].childNodes[0].nodeValue
                  outline = symbol.getElementsByTagName("FillColor")[0].childNodes[0].nodeValue
                  size = symbol.getElementsByTagName("FillColor")[0].childNodes[0].nodeValue
                  printMessage("FillColor->R,G,B: " + str(fill))
                  printMessage("OutlineColor->R,G,B: " + str(outline))
                  printMessage("Size: "+str(size))
               else:
                  color = symbol.getElementsByTagName("Color")[0].childNodes[0].nodeValue
                  printMessage("Color: " + str(color))
           else:
               printMessage("Not found: " + lyr.name)
               return "",""
        else:
           return "",""
        return lyr.name,drawingInfo


#def getFeatures(lyr):
   #loop through each row in layer

   #for row in arcpy.SearchCursor(lyr):
    #df.extent = row.SHAPE.extent #Set the dataframe extent to the extent of the feature
    #df.scale = df.scale * 1.07 #Optionally give the shape a bit of padding around the edges

   #features = {

   #}
   #return features


#   if 'layerDefinition' in layer:
#      if 'fields' in layer['layerDefinition']:
#              for field in layer['layerDefinition']['fields']:
#                  fieldInfos = None
#                  if field['type'] == 'esriFieldTypeOID':
#                      oidFldName = field['name']
#                      fieldInfos = {
#                          'fieldName':field['name'],
#                          'label':field['alias'],
#                          'isEditable':False,
#                          'tooltip':'',
#                          'visible':False,
#                          'format':None,
#                          'stringFieldOption':'textbox'
#                      }
#
#                  elif field['type'] == 'esriFieldTypeInteger':
#                      fieldInfos = {
#                          'fieldName':field['name'],
#                          'label':field['alias'],
#                          'isEditable':True,
#                          'tooltip':'',
#                          'visible':True,
#                          'format':{
#                              'places':0,
#                              'digitSeparator':True
#                          },
#                          'stringFieldOption':'textbox'
#                      }
#                  elif field['type'] == 'esriFieldTypeDouble':
#                      fieldInfos = {
#                          'fieldName':field['name'],
#                          'label':field['alias'],
#                          'isEditable':True,
#                          'tooltip':'',
#                          'visible':True,
#                          'format':{
#                              'places':2,
#                              'digitSeparator':True
#                              },
#                          'stringFieldOption':'textbox'
#                      }
#                  elif field['type'] == 'esriFieldTypeString':
#                      fieldInfos = {
#                          'fieldName':field['name'],
#                          'label':field['alias'],
#                          'isEditable':True,
#                          'tooltip':'',
#                          'visible':True,
#                          'format':None,
#                          'stringFieldOption':'textbox'
#                      }
#                  else:
#                      fieldInfos = {
#                          'fieldName':field['name'],
#                          'label':field['alias'],
#                          'isEditable':True,
#                          'tooltip':'',
#                          'visible':True,
#                          'format':None,
#                          'stringFieldOption':'textbox'
#                      }
#                  if fieldInfos is not None:
#                      popInfo.append(fieldInfos)

#def getFieldInfos(layer):
#    # Create a describe object
#    desc = arcpy.Describe(layer)
#
#    # If a feature layer, continue
#    if desc.dataType == "FeatureLayer":
#
#        # Create a fieldinfo object
#        field_info = desc.fieldInfo
#
#        # Use the count property to iterate through all the fields
#        for index in xrange(0, field_info.count):
#            # Print fieldinfo properties
#            printMessage("Field Name: {0}".format(field_info.getFieldName(index)))
#            printMessage("\tNew Name:   {0}".format(field_info.getNewName(index)))
#            printMessage("\tSplit Rule: {0}".format(field_info.getSplitRule(index)))
#            printMessage("\tVisible:    {0}".format(field_info.getVisible(index)))




def printMessage(str):
  if sys.executable.find("python.exe") != -1:
     print(str)
  else:
     try:
       arcpy.AddMessage(str)
     except Exception as e:
       print(str)


def main():
    tbx=Toolbox()
    tool=CreateNewProject()
    #tool.execute(tool.getParameterInfo(),r"C:\hpl\distribution\aar\leasecompliance2014.gdb.mxd")
    #mxd,server,user,outputfolder
    #tool.execute(tool.getParameterInfo(),r"C:\Users\steve\Documents\ArcGIS\Packages\leasecompliance2016_B4A776C0-3F50-4B7C-ABEE-76C757E356C7\v103\leasecompliance2016.mxd|gis.biz.tm|shale|D:\workspace\go\src\github.com\traderboy\arcrestgo\leasecompliance2016")
    #tool.execute(tool.getParameterInfo(),r"C:\Users\steve\Documents\ArcGIS\Packages\leasecompliance2016_B4A776C0-3F50-4B7C-ABEE-76C757E356C7\v103\leasecompliance2016.mxd|reais.x10host.com|shale|D:\workspace\go\src\github.com\traderboy\arcrestgo\leasecompliance2016")
    
    tool.execute(tool.getParameterInfo(),r"C:\Users\steve\Documents\ArcGIS\Packages\leasecompliance2016_B629916B-D98A-42C5-B9E1-336B123CECDF\v103\leasecompliance2016.mxd|reais.x10host.com|shale|C:\docker\src\github.com\traderboy\arcrestgo\leasecompliance2016")
    #tool.execute(tool.getParameterInfo(),r"D:\workspace\hpl\distribution\aar\Accommodation Agreement Rentals.mxd")
    #arcpy.ImportToolbox ("C:/Users/steve/git/arcservice/Createarcgisprojecttool.pyt")
    #arcpy.arcservices.CreateNewProject()

if __name__ == '__main__':
    if sys.executable.find("python.exe") != -1:
       main()