import json
import unittest
from packaging import ManifestInfo, OfferInfo, DeploymentType


class TestManifest(unittest.TestCase):
    def test_deployment_type_enum(self):
        self.assertEqual(DeploymentType.terraform.value, "terraform")
        self.assertEqual(DeploymentType.bicep.value, "bicep")

    def test_manifest_info_construction(self):
        manifest = ManifestInfo(main_template="main.tf", deployment_type=DeploymentType.terraform)

        self.assertEqual(manifest.main_template, "main.tf")
        self.assertEqual(manifest.deployment_type, DeploymentType.terraform)

        self.assertIsInstance(manifest.offer, OfferInfo)

    def test_manifest_info_serialization(self):
        manifest = ManifestInfo(main_template="main.tf", deployment_type=DeploymentType.terraform)

        json = manifest.serialize()
        self.assertEqual(json["mainTemplate"], manifest.main_template)
        self.assertEqual(json["deploymentType"], manifest.deployment_type.value)
    
    def test_manifest_info_to_json(self):
        manifest = ManifestInfo(main_template="main.tf", deployment_type=DeploymentType.terraform)

        json_str = manifest.to_json()
        from_json = ManifestInfo().deserialize(json.loads(json_str))
        self.assertEqual(from_json.main_template, manifest.main_template)

    def test_manifest_info_serialization(self):
        manifest = ManifestInfo(main_template="main.tf", deployment_type=DeploymentType.terraform)


    def test_offer_info(self):
        offer = OfferInfo(name = "test", description = "test")
        self.assertEqual(offer.name, "test")
        self.assertEqual(offer.name, "test")

        offer = OfferInfo()
        self.assertEqual(offer.name, "")
        self.assertEqual(offer.name, "")


if __name__ == '__main__':
    unittest.main()