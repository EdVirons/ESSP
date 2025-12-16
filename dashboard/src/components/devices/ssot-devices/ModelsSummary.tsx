import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';

interface DeviceModel {
  model: string;
  count: number;
}

interface ModelsSummaryProps {
  models: DeviceModel[];
  makes: string[];
}

export function ModelsSummary({ models, makes }: ModelsSummaryProps) {
  if (models.length === 0) {
    return null;
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">Device Models</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex flex-wrap gap-2">
          {models.map((model) => (
            <Badge key={model.model} variant="outline" className="text-sm">
              {model.model} ({model.count})
            </Badge>
          ))}
        </div>
        {makes.length > 0 && (
          <div className="mt-4 pt-4 border-t">
            <p className="text-sm text-gray-500 mb-2">Makes:</p>
            <div className="flex flex-wrap gap-2">
              {makes.map((make) => (
                <span key={make} className="text-sm text-gray-700 bg-gray-100 px-2 py-1 rounded">
                  {make}
                </span>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
