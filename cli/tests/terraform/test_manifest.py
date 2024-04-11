import json
import unittest
from unittest.mock import patch
from modm.installer import ManifestInfo, OfferProperties, DeploymentType


class TestManifest(unittest.TestCase):
    def test_deployment_type_enum(self):
        self.assertEqual(DeploymentType.terraform.value, "terraform")
        self.assertEqual(DeploymentType.arm.value, "arm")

    def test_manifest_info_construction(self):
        manifest = ManifestInfo(solution_template="main.tf", deployment_type=DeploymentType.terraform)

        self.assertEqual(str(manifest.solution_template), "main.tf")
        self.assertEqual(manifest.deployment_type, DeploymentType.terraform)

        self.assertIsInstance(manifest.offer, OfferProperties)

    @patch('modm.installer.ManifestInfo._compile_bicep_template')
    def test_manifest_info_uses_file_ext_to_set_deployment_type(self, mock_compile_bicep_template):
        mock_compile_bicep_template.return_value = None
        
        manifest = ManifestInfo(solution_template="main.tf")
        self.assertEqual(manifest.deployment_type, DeploymentType.terraform)

        manifest = ManifestInfo(solution_template="main.bicep")
        self.assertEqual(manifest.deployment_type, DeploymentType.arm)

    def test_manifest_info_serialization(self):
        manifest = ManifestInfo(solution_template="main.tf", deployment_type=DeploymentType.terraform)

        json = manifest.serialize()
        self.assertEqual(json["mainTemplate"], manifest.solution_template)
        self.assertEqual(json["deploymentType"], manifest.deployment_type.value)
    
    def test_manifest_info_to_json(self):
        manifest = ManifestInfo(solution_template="main.tf", deployment_type=DeploymentType.terraform)

        json_str = manifest.to_json()
        from_json: ManifestInfo = manifest.deserialize(json.loads(json_str))
        self.assertEqual(from_json.solution_template, manifest.solution_template)

    def test_manifest_info_serialization(self):
        manifest = ManifestInfo(solution_template="main.tf", deployment_type=DeploymentType.terraform)


    def test_offer_info(self):
        offer = OfferProperties(name = "test", description = "test")
        self.assertEqual(offer.name, "test")
        self.assertEqual(offer.name, "test")

        offer = OfferProperties()
        self.assertEqual(offer.name, "")
        self.assertEqual(offer.name, "")
